package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	model "fgl-backend/app/model"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// GetAnnouncement returns an announcement
func GetAnnouncement(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)
	id := vars["aid"]
	ann := getAnnOr404(db, id, w, r)
	if ann == nil {
		return
	}
	RespondJSON(w, http.StatusOK, ann)
}

// GetAnnouncements returns all announcements
func GetAnnouncements(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	announcements := []model.Announcement{}
	db.Find(&announcements)
	RespondJSON(w, http.StatusOK, announcements)
}

// GetAnnouncementsFromUser gets all announcements from a specific user
func GetAnnouncementsFromUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	anns := []model.Announcement{}

	if err := db.Find(&anns, model.Announcement{Author: user}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, anns)
}

// MakeAnnouncement creates a new announcement
func MakeAnnouncement(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	ann := model.Announcement{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ann); err != nil {
		fmt.Println("error decoding announcement:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&ann).Error; err != nil {
		fmt.Println("error saving announcement:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, ann)
	fmt.Println("New Announcement By " + ann.Author + ": " + ann.Title + strconv.Itoa(int(ann.AID)))
}

// UpdateAnnouncement updates an announcement
func UpdateAnnouncement(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["aid"]
	ann := getAnnOr404(db, id, w, r)
	if ann == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ann); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	if err := db.Save(&ann).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, ann)
}

// DeleteAnnouncement deletes an announcement
func DeleteAnnouncement(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["aid"]

	ann := getAnnOr404(db, id, w, r)
	if ann == nil {
		return
	}
	if err := db.Delete(&ann).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

func getAnnOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Announcement {
	ann := model.Announcement{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&ann, model.Announcement{AID: idInt}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &ann
}
