package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "fgl-backend/app/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// GetUser returns a user
func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	user := getUserOr404(db, username, w, r)
	if user == nil {
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// GetAllUsers returns all users
func GetAllUsers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	db.Find(&users)
	respondJSON(w, http.StatusOK, users)
}

// UserLogin handles user login
func UserLogin(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	credentials := model.Credentials{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&credentials); err != nil {
		fmt.Println("error decoding credentials:", err)
	}

	user := getUserOr404(db, credentials.Username, w, r)
	if user == nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		fmt.Println("error", err)
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims := &model.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_key"))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		fmt.Println("error generating token")
		return
	}

	respondJSON(w, http.StatusOK, tokenString)
	fmt.Println("user has logged in:", user.Username)
}

// RegisterUser handles user registration
func RegisterUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		fmt.Println("error decoding user information:", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		fmt.Println("error hashing password:", err)
	}
	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		fmt.Println("error saving user to database:", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, user)
	fmt.Println("Registered new user:", user.Username)
}

// UpdateUser updates user information for specified user
func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username := vars["username"]
	user := getUserOr404(db, username, w, r)
	if user == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// DeleteUser deletes specified user
func DeleteUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username := vars["username"]
	user := getUserOr404(db, username, w, r)
	if user == nil {
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func getUserOr404(db *gorm.DB, username string, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{Username: username}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}
