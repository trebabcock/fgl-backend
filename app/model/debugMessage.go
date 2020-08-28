package model

import (
	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DebugMessage holds the debug message parameters
type DebugMessage struct {
	gorm.Model
	Author string `json:"author"`
	Body   string `json:"body"`
	Time   string `json:"time"`
}
