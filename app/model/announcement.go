package model

import (
	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Announcement holds the announcement parameters
type Announcement struct {
	gorm.Model
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	AID    int64  `gorm:"AUTO_INCREMENT" json:"aid"`
}
