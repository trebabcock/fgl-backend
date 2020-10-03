package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DebugMessage struct {
	gorm.Model
	Author string `json:"author"`
	Body   string `json:"body"`
	Time   string `json:"time"`
}
