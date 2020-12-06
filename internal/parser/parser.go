package parser

import (
	"github.com/pkg/errors"
)

type Result struct {
	Tracker  string
	Name     string
	SubName  string
	LinkSite string
	LinkDL   string
	Year     int
	Category Category
	Tags     []string
	Formats  []string
}

type Category string

const (
	tv    Category = "tv"
	movie Category = "movie"
	music Category = "music"
	app   Category = "app"
	game  Category = "game"
	ebook Category = "ebook"
	xxx   Category = "xxx"
	anime Category = "anime"
)

var (
	ErrCannotParse   = errors.New("Failed to parse message")
	ErrInvalidParser = errors.New("Invalid/Unknown IParser")
)
