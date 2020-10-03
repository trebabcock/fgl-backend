package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type LabReport struct {
	gorm.Model
	Author  string `json:"author"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	RID     int64  `gorm:"AUTO_INCREMENT" json:"rid"`
}
