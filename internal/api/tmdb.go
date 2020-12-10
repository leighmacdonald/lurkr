package api

import (
	"github.com/cyruzin/golang-tmdb"
	"github.com/pkg/errors"
)

var (
	tmdbClient    *tmdb.Client
	ErrNotFound   = errors.New("Results not found")
	ErrInvalidKey = errors.New("Invalid api key")
)

func SetupTMDB(key string) error {
	if key == "" {
		return ErrInvalidKey
	}
	client, err := tmdb.Init(key)
	if err != nil {
		return err
	}
	tmdbClient = client
	return nil
}

type ScoreFn func(string) (float64, error)

func TMDBScore(name string) (float64, error) {
	if tmdbClient == nil {
		return 0, ErrInvalidKey
	}
	return 0, nil
}
