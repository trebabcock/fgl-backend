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

func GetLabReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["rid"]
	lab := getLabOr404(db, id, w, r)
	if lab == nil {
		return
	}
	respondJSON(w, http.StatusOK, lab)
}

func GetLabReports(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	reports := []model.LabReport{}
	db.Find(&reports)
	respondJSON(w, http.StatusOK, reports)
}

func MakeLabReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	lab := model.LabReport{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lab); err != nil {
		fmt.Println("error decoding report:", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&lab).Error; err != nil {
		fmt.Println("error saving report:", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, lab)
	fmt.Println("New Lab Report By " + lab.Author + ": " + lab.Title)
}

func getLabOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.LabReport {
	lab := model.LabReport{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&lab, model.LabReport{RID: idInt}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &lab
}
