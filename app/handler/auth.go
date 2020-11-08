package handler

import (
	"crypto/sha256"
	"encoding/binary"
	"fgl-backend/app/model"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// Authorize checks if a given auth code is authentic
func Authorize(db *gorm.DB, w http.ResponseWriter, r *http.Request) bool {

	r.ParseForm()
	username := r.Form.Get("username")
	authcode := r.Form.Get("authcode")

	user := model.User{}
	if err := db.First(&user, model.User{Username: username}).Error; err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error())
		log.Println(username, authcode, err.Error())
		return false
	}
	return user.AuthCode == authcode
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

// NewAuthCode creates a new auth code
func NewAuthCode(user *model.User) string {
	codehash := sha256.New()
	codehash.Write([]byte(user.Username))
	codehash.Write([]byte(user.Password))
	timeslice := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeslice, uint64(time.Now().UnixNano()))
	codehash.Write(timeslice)
	return string(codehash.Sum(nil))
}
