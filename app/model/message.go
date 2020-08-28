package model

import (
	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Message holds the message parameters
type Message struct {
	gorm.Model
	Author string `json:"author"`
	Body   string `json:"body"`
	Time   string `json:"time"`
}
