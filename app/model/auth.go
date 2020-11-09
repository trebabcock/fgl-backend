package model

import "github.com/jinzhu/gorm"

// Authenticator is a download authenticator
type Authenticator struct {
	gorm.Model
	Code string `json:"auth_code"`
}
