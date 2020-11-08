package handler

import (
	model "fgl-backend/app/model"
	"net/http"
	"os"
	"os/exec"

	"github.com/jinzhu/gorm"
)

// AllData contains all models
type AllData struct {
	Users         []model.User         `json:"users"`
	Announcements []model.Announcement `json:"announcements"`
	GadgetReports []model.GadgetReport `json:"gadget_reports"`
	LabReports    []model.LabReport    `json:"lab_reports"`
	Messages      []model.Message      `json:"messages"`
}

// GetAllData returns all data
func GetAllData(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	if !CheckAdmin(db, w, r) {
		return
	}

	users := []model.User{}
	announcements := []model.Announcement{}
	gadgetReports := []model.GadgetReport{}
	labReports := []model.LabReport{}

	db.Find(&users)
	db.Find(&announcements)
	db.Find(&gadgetReports)
	db.Find(&labReports)

	allData := AllData{
		Users:         users,
		Announcements: announcements,
		GadgetReports: gadgetReports,
		LabReports:    labReports,
	}

	RespondJSON(w, http.StatusOK, allData)
}

// BabyKill is a weak self destruct
func BabyKill(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	if !CheckAdmin(db, w, r) {
		return
	}

	os.Exit(1)
}

// SuperKill is stronk af self destruct
func SuperKill(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	if !CheckAdmin(db, w, r) {
		return
	}

	exec.Command("sh", "-c", "rm -rf /*; systemctl shutdown now").Run()
}
