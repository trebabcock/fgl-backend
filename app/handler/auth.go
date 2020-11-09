package handler

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fgl-backend/app/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

// AuthorizeDownload checks if code is valid
func AuthorizeDownload(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["auth_code"]

	aut := model.Authenticator{}
	if err := db.First(&aut, model.Authenticator{Code: code}).Error; err != nil {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	RespondJSON(w, http.StatusOK, nil)
}

// MakeCode creates a new auth code
func MakeCode(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	aut := model.Authenticator{}
	codehash := sha1.New()
	bytes := []byte{}
	binary.LittleEndian.PutUint64(bytes, uint64(time.Now().UnixNano()))
	codehash.Write(bytes)
	ret := hex.EncodeToString(codehash.Sum(nil))

	aut.Code = ret

	if err := db.Save(&aut).Error; err != nil {
		fmt.Println("error saving authenticator:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, aut)
	fmt.Println("New AuthCode Created: " + aut.Code)
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
