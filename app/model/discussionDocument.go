package model

import (
	"github.com/jinzhu/gorm"
	// Dialect for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DiscussionDocument is a discussion document
type DiscussionDocument struct {
	gorm.Model
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	DID    int64  `gorm:"AUTO_INCREMENT" json:"did"`
}
