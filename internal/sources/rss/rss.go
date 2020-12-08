package rss

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"io"
)

var (
	ErrParse = errors.New("Failed to parse RSS/Atom feed")
	sources  = map[string]Feed{}
)

type Feed struct {
	cfg    config.TrackerConfig
	parser *gofeed.Parser
}

func (s Feed) Parse(reader io.Reader) (*gofeed.Feed, error) {
	feed, err := s.parser.Parse(reader)
	if err != nil {
		return nil, errors.Wrapf(ErrParse, err.Error())
	}
	return feed, nil
}

func New(cfg config.TrackerConfig) *Feed {
	parser := gofeed.NewParser()
	return &Feed{cfg: cfg, parser: parser}
}
