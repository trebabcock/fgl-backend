package handler

import (
	"fgl-backend/app/model"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Authorize checks if a given auth code is authentic
func Authorize(db *gorm.DB, w http.ResponseWriter, r *http.Request) bool {

	r.ParseForm()
	username := r.Form.Get("username")

	user := model.User{}
	if err := db.First(&user, model.User{Username: username}).Error; err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error())
		return false
	}
	return true
}

// CheckAdmin checks if a user has admin status
func CheckAdmin(db *gorm.DB, w http.ResponseWriter, r *http.Request) bool {

	r.ParseForm()
	username := r.Form.Get("username")

	user := model.User{}
	if err := db.First(&user, model.User{Username: username}).Error; err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error())
		log.Println(username, err.Error())
		return false
	}
	return user.Admin
}
