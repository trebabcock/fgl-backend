package handler

import (
	"encoding/json"
	model "fgl-backend/app/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetDiscussionDocument gets a discussion report
func GetDiscussionDocument(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)
	id := vars["rid"]
	discussion := getDiscussionOr404(db, id, w, r)
	if discussion == nil {
		return
	}
	RespondJSON(w, http.StatusOK, discussion)
}

// GetDiscussionDocuments gets all discussion reports
func GetDiscussionDocuments(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	reports := []model.DiscussionDocument{}
	db.Find(&reports)
	RespondJSON(w, http.StatusOK, reports)
}

// GetDiscussionsFromUser gets all announcements from a specific user
func GetDiscussionsFromUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	dis := []model.DiscussionDocument{}

	if err := db.Find(&dis, model.DiscussionDocument{Author: user}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, dis)
}

// MakeDiscussionDocument creates a new discussion report
func MakeDiscussionDocument(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	discussion := model.DiscussionDocument{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&discussion); err != nil {
		fmt.Println("error decoding report:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&discussion).Error; err != nil {
		fmt.Println("error saving report:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, discussion)
	fmt.Println("New Discussion Report By " + discussion.Author + ": " + discussion.Title)
}

// UpdateDiscussionDocument updates an announcement
func UpdateDiscussionDocument(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]
	lr := getDiscussionOr404(db, id, w, r)
	if lr == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lr); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	if err := db.Save(&lr).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, lr)
}

// DeleteDiscussionDocument deletes an announcement
func DeleteDiscussionDocument(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]

	lr := getDiscussionOr404(db, id, w, r)
	if lr == nil {
		return
	}
	if err := db.Delete(&lr).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

func getDiscussionOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.DiscussionDocument {
	discussion := model.DiscussionDocument{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&discussion, model.DiscussionDocument{DID: idInt}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &discussion
}
