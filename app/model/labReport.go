package model

import (
	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// LabReport holds the lab report parameters
type LabReport struct {
	gorm.Model
	Author  string `json:"author"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	RID     int64  `gorm:"AUTO_INCREMENT" json:"rid"`
}
