package model

import (
	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GadgetReport holds the gadget report parameters
type GadgetReport struct {
	gorm.Model
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	RID    int64  `gorm:"AUTO_INCREMENT" json:"aid"`
}
