package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	// Postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User holds the user parameters
type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// Credentials holds user credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims holds jwt claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
