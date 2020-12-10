package internal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Release struct {
	gorm.Model
	Hash      string
	Tracker   string
	CreatedOn time.Time
}

func NewDb(dsn string) (*gorm.DB, error) {
	newDb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := newDb.AutoMigrate(&Release{}); err != nil {
		return nil, err
	}
	return newDb, nil
}

func GetRelease(ih string) (*Release, error) {
	var rls Release
	err := db.First(&rls, "hash = ?", ih).Error
	return &rls, err
}

func InsertRelease(hash string, tracker string) error {
	return db.Create(&Release{Hash: hash, Tracker: tracker, CreatedOn: time.Now()}).Error
}
