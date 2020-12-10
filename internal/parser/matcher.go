package parser

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	log "github.com/sirupsen/logrus"
)

func Match(cfg config.TrackerConfig, result *Result) bool {
	if ContainsAnyStrings(cfg.Filters.TagsExcluded, result.Tags) {
		log.Debugf("Skipped release due to exluded tags: %v", result.Name)
		return false
	}
	if len(cfg.Filters.TagsAllowed) > 0 && !ContainsAnyStrings(cfg.Filters.TagsAllowed, result.Tags) {
		log.Debugf("Skipped release due to no matching tag: %v", result.Name)
		return false
	}
	return true
}

func ContainsString(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func ContainsAnyStrings(needles []string, haystack []string) bool {
	for _, needle := range needles {
		if ContainsString(needle, haystack) {
			return true
		}
	}
	return false
}
