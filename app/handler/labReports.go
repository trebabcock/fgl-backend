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

// GetLabReport gets a lab report
func GetLabReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)
	id := vars["rid"]
	lab := getLabOr404(db, id, w, r)
	if lab == nil {
		return
	}
	RespondJSON(w, http.StatusOK, lab)
}

// GetLabReports gets all lab reports
func GetLabReports(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	reports := []model.LabReport{}
	db.Find(&reports)
	RespondJSON(w, http.StatusOK, reports)
}

// MakeLabReport creates a new lab report
func MakeLabReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	lab := model.LabReport{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lab); err != nil {
		fmt.Println("error decoding report:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&lab).Error; err != nil {
		fmt.Println("error saving report:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, lab)
	fmt.Println("New Lab Report By " + lab.Author + ": " + lab.Title)
}

// UpdateLabReport updates an announcement
func UpdateLabReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]
	lr := getLabOr404(db, id, w, r)
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

// DeleteLabReport deletes an announcement
func DeleteLabReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]

	lr := getLabOr404(db, id, w, r)
	if lr == nil {
		return
	}
	if err := db.Delete(&lr).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

func getLabOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.LabReport {
	lab := model.LabReport{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&lab, model.LabReport{RID: idInt}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &lab
}
