package model

import (
	"github.com/jinzhu/gorm"
	// Dialect for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Vote is a vote
type Vote struct {
	gorm.Model
	Caller   string   `json:"caller"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Type     string   `json:"type"`
	Yeas     []string `gorm:"many2many" json:"yeas"`
	Nays     []string `gorm:"many2many" json:"nays"`
	Complete bool     `json:"complete"`
	Passed   bool     `json:"passed"`
	VID      int64    `gorm:"AUTO_INCREMENT" json:"vid"`
}

// VoterReport is a voter report
type VoterReport struct {
	gorm.Model
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	VRID   int64  `gorm:"AUTO_INCREMENT" json:"vrid"`
}
