package model

import (
	"github.com/jinzhu/gorm"
	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DBMigrate migrates the postgres database
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{}, &Announcement{}, &LabReport{}, &GadgetReport{}, &DebugMessage{})
	return db
}
