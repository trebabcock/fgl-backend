package model

import (
	"github.com/jinzhu/gorm"
	// Dialect for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ProjectSubmission is a project submission
type ProjectSubmission struct {
	gorm.Model
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	PID    int64  `gorm:"AUTO_INCREMENT" json:"pid"`
}

// Project is a project
type Project struct {
	gorm.Model
	Name       string   `json:"name"`
	Lead       string   `json:"lead"`
	Assistants []string `json:"assistants"`
	Updates    []string `json:"updates"`
}
