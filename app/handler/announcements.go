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
	vars := mux.Vars(r)
	id := vars["aid"]
	ann := getAnnOr404(db, id, w, r)
	if ann == nil {
		return
	}
	respondJSON(w, http.StatusOK, ann)
}

// GetAnnouncements returns all announcements
func GetAnnouncements(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	announcements := []model.Announcement{}
	db.Find(&announcements)
	respondJSON(w, http.StatusOK, announcements)
}

// MakeAnnouncement creates a new announcement
func MakeAnnouncement(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ann := model.Announcement{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ann); err != nil {
		fmt.Println("error decoding announcement:", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&ann).Error; err != nil {
		fmt.Println("error saving announcement:", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, ann)
	fmt.Println("New Announcement By " + ann.Author + ": " + ann.Title + strconv.Itoa(int(ann.AID)))
}

func getAnnOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Announcement {
	ann := model.Announcement{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&ann, model.Announcement{AID: idInt}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &ann
}
