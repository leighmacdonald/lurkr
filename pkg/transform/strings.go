package transform

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	rxFmt             *regexp.Regexp
	ErrFmtMissingArgs = errors.New("Missing required formatting arguments")
)

func ToInt(value string) int {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1
	}
	return int(v)
}

func TrimStrings(stringSlice []string) []string {
	var trimmed []string
	for _, str := range stringSlice {
		trimmed = append(trimmed, strings.TrimSpace(str))
	}
	return trimmed
}

func LowerStrings(stringSlice []string) []string {
	var lowered []string
	for _, str := range stringSlice {
		lowered = append(lowered, strings.ToLower(str))
	}
	return lowered
}

func NormalizeStrings(stringSlice []string) []string {
	return LowerStrings(TrimStrings(stringSlice))
}

func FindYear(release string) int {
	pcs := strings.Split(strings.Replace(release, ".", " ", -1), " ")
	for _, p := range pcs {
		year := ToInt(p)
		if year > 1900 {
			return year
		}
	}
	return 0
}

type FArgs map[string]interface{}

// Format will do string formatting, replacing placeholders with the supplied values.
// If the input format does not have all its placeholders replaced, it will return
// an ErrFmtMissingArgs error. Extra FArgs that do not have a matching placeholder
// in the format string are not considered an error and are safely ignored.
//
// Format("a {placeholder}, FArgs("placeholder": "replacement"))
func Format(format string, args FArgs) (string, error) {
	matches := rxFmt.FindAllStringSubmatch(format, -1)
	var (
		found   []string
		missing []string
	)
	for _, ma := range matches {
		found = append(found, ma[1])
	}
	for _, k := range found {
		isHere := false
		for f := range args {
			if f == k {
				isHere = true
				break
			}
		}
		if !isHere {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return "", errors.Wrapf(ErrFmtMissingArgs, "Missing arguments: %s", strings.Join(missing, ", "))
	}
	var replacements []string
	for k, v := range args {
		replacements = append(replacements, "{"+k+"}", fmt.Sprint(v))
	}
	return strings.NewReplacer(replacements...).Replace(format), nil
}

func init() {
	rxFmt = regexp.MustCompile(`{(.+?)}`)
}
