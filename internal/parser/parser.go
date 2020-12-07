package parser

import (
	"github.com/pkg/errors"
	"strings"
)

type Result struct {
	Tracker    string
	Name       string
	SubName    string
	LinkSite   string
	LinkDL     string
	Year       int
	Episode    int
	Season     int
	Group      string
	Resolution string
	Source     string
	Codec      string
	// Full Seasons = pack
	Pack     bool
	Category Category
	Tags     []string
	Formats  []string
}

type Category string

const (
	TV          Category = "tv"
	Movie       Category = "movie"
	Music       Category = "music"
	AppPC       Category = "app_pc"
	AppMac      Category = "app_mac"
	GamePC      Category = "game_pc"
	GameConsole Category = "game_console"
	Ebook       Category = "ebook"
	Mobile      Category = "mobile"
	XXX         Category = "xxx"
	Anime       Category = "anime"
	Unknown     Category = "unknown"
)

var (
	ErrCannotParse   = errors.New("Failed to parse message")
	ErrInvalidParser = errors.New("Invalid/Unknown Parser")
)

func NewResult(driver string) *Result {
	return &Result{
		Tracker: driver,
		Formats: []string{},
		Tags:    []string{},
	}
}

func FindCategory(name string, mapping CategoryMap) Category {
	if mapping == nil {
		mapping = CatMap
	}
	name = strings.ToLower(name)
	for cat, values := range mapping {
		for _, value := range values {
			if value == name {
				return cat
			}
		}
	}
	return Unknown
}

type CategoryMap map[Category][]string

var (
	// Common mappings to categories. Many trackers will require custom mappings ontop of these
	CatMap = CategoryMap{
		Movie:       {"movie", "movies"},
		TV:          {"tv", "television"},
		GamePC:      {"game", "games"},
		GameConsole: {"ps3", "ps4", "ps5", "wii", "xbox360", "nds", "psp"},
		Ebook:       {"ebook", "e-book", "ebooks", "e-books"},
		Music:       {"music", "flac", "mp3"},
		AppPC:       {"app", "application"},
		AppMac:      {"apple", "mac", "macos"},
		Mobile:      {"android", "ios"},
		Anime:       {"anime"},
		XXX:         {"xxx"},
	}
)
