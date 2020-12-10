package extern_rules

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ImportIrssiRules(directory string) ([]IrssiRule, error) {
	rxExample := regexp.MustCompile(`<!--\s?(.+?)\s?-->`)
	var rules []IrssiRule
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".tracker") {
			return nil
		}
		fullPath := filepath.Join(directory, info.Name())
		var rule IrssiRule
		b, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return err
		}
		if err := xml.NewDecoder(bytes.NewReader(b)).Decode(&rule); err != nil {
			return err
		}
		p := strings.Split(string(b), "END LICENSE BLOCK")
		lines := strings.ReplaceAll(strings.ReplaceAll(p[1], "\r\n", "\n"), "\n", "")
		var exampleLines []string
		eMatch := rxExample.FindAllString(lines, -1)
		for _, match := range eMatch {
			exampleLines = append(exampleLines, match)
		}
		if len(exampleLines) == 0 {
			// Skip rules w/o example data
			return nil
		}
		rule.ExampleLines = exampleLines
		rules = append(rules, rule)
		return nil
	})
	return rules, err
}

type IrssiRule struct {
	ExampleLines []string
	XMLName      xml.Name `xml:"trackerinfo"`
	Text         string   `xml:",chardata"`
	Type         string   `xml:"type,attr"`
	ShortName    string   `xml:"shortName,attr"`
	LongName     string   `xml:"longName,attr"`
	SiteName     string   `xml:"siteName,attr"`
	Settings     struct {
		Text               string `xml:",chardata"`
		GazelleDescription string `xml:"gazelle_description"`
		GazelleAuthkey     struct {
			Text       string `xml:",chardata"`
			PasteGroup string `xml:"pasteGroup,attr"`
			PasteRegex string `xml:"pasteRegex,attr"`
		} `xml:"gazelle_authkey"`
		GazelleTorrentPass string `xml:"gazelle_torrent_pass"`
	} `xml:"settings"`
	Servers struct {
		Text   string `xml:",chardata"`
		Server struct {
			Text           string `xml:",chardata"`
			Network        string `xml:"network,attr"`
			ServerNames    string `xml:"serverNames,attr"`
			ChannelNames   string `xml:"channelNames,attr"`
			AnnouncerNames string `xml:"announcerNames,attr"`
		} `xml:"server"`
	} `xml:"servers"`
	Parseinfo struct {
		Text         string `xml:",chardata"`
		Linepatterns struct {
			Text    string `xml:",chardata"`
			Extract struct {
				Text  string `xml:",chardata"`
				Regex struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"regex"`
				Vars struct {
					Text string `xml:",chardata"`
					Var  []struct {
						Text string `xml:",chardata"`
						Name string `xml:"name,attr"`
					} `xml:"var"`
				} `xml:"vars"`
			} `xml:"extract"`
		} `xml:"linepatterns"`

		Multilinepatterns struct {
			Text    string `xml:",chardata"`
			Extract []struct {
				Text  string `xml:",chardata"`
				Regex struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"regex"`
				Vars struct {
					Text string `xml:",chardata"`
					Var  []struct {
						Text string `xml:",chardata"`
						Name string `xml:"name,attr"`
					} `xml:"var"`
				} `xml:"vars"`
			} `xml:"extract"`
		} `xml:"multilinepatterns"`
		Linematched struct {
			Text string `xml:",chardata"`
			Var  []struct {
				Text   string `xml:",chardata"`
				Name   string `xml:"name,attr"`
				String []struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"string"`
				Var []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"var"`
			} `xml:"var"`
			If []struct {
				Text   string `xml:",chardata"`
				Srcvar string `xml:"srcvar,attr"`
				Regex  string `xml:"regex,attr"`
				Var    struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					String []struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"string"`
					Var []struct {
						Text string `xml:",chardata"`
						Name string `xml:"name,attr"`
					} `xml:"var"`
				} `xml:"var"`
			} `xml:"if"`
			Extract []struct {
				Text     string `xml:",chardata"`
				Srcvar   string `xml:"srcvar,attr"`
				Optional string `xml:"optional,attr"`
				Regex    struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"regex"`
				Vars struct {
					Text string `xml:",chardata"`
					Var  []struct {
						Text string `xml:",chardata"`
						Name string `xml:"name,attr"`
					} `xml:"var"`
				} `xml:"vars"`
			} `xml:"extract"`
			Varreplace struct {
				Text    string `xml:",chardata"`
				Name    string `xml:"name,attr"`
				Srcvar  string `xml:"srcvar,attr"`
				Regex   string `xml:"regex,attr"`
				Replace string `xml:"replace,attr"`
			} `xml:"varreplace"`
			Setregex struct {
				Text     string `xml:",chardata"`
				Srcvar   string `xml:"srcvar,attr"`
				Regex    string `xml:"regex,attr"`
				VarName  string `xml:"varName,attr"`
				NewValue string `xml:"newValue,attr"`
			} `xml:"setregex"`
			Extracttags struct {
				Text     string `xml:",chardata"`
				Srcvar   string `xml:"srcvar,attr"`
				Split    string `xml:"split,attr"`
				Setvarif []struct {
					Text    string `xml:",chardata"`
					VarName string `xml:"varName,attr"`
					Regex   string `xml:"regex,attr"`
				} `xml:"setvarif"`
				Regex struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"regex"`
			} `xml:"extracttags"`
		} `xml:"linematched"`
		Ignore struct {
			Text  string `xml:",chardata"`
			Regex []struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value,attr"`
			} `xml:"regex"`
		} `xml:"ignore"`
	} `xml:"parseinfo"`
}
