package model

import (
	"github.com/jinzhu/gorm"
	// Dialect for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GadgetReport is a gadget report
type GadgetReport struct {
	gorm.Model
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	RID    int64  `gorm:"AUTO_INCREMENT" json:"rid"`
}
