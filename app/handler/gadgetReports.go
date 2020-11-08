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

// GetGadgetReport returns a single gadget report
func GetGadgetReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)
	id := vars["rid"]
	gadget := getGadgetOr404(db, id, w, r)
	if gadget == nil {
		return
	}
	RespondJSON(w, http.StatusOK, gadget)
}

// GetGadgetReports returns all gadget reports
func GetGadgetReports(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	reports := []model.GadgetReport{}
	db.Find(&reports)
	RespondJSON(w, http.StatusOK, reports)
}

// MakeGadgetReport makes a new gadget report
func MakeGadgetReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	gadget := model.GadgetReport{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gadget); err != nil {
		fmt.Println("error decoding report:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&gadget).Error; err != nil {
		fmt.Println("error saving report:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, gadget)
	fmt.Println("New Gadget Report By " + gadget.Author + ": " + gadget.Title)
}

// UpdateGadgetReport updates an announcement
func UpdateGadgetReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]
	gr := getGadgetOr404(db, id, w, r)
	if gr == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&gr); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
	}
	defer r.Body.Close()

	if err := db.Save(&gr).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, gr)
}

// DeleteGadgetReport deletes an announcement
func DeleteGadgetReport(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]

	gr := getGadgetOr404(db, id, w, r)
	if gr == nil {
		return
	}
	if err := db.Delete(&gr).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

func getGadgetOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.GadgetReport {
	gadget := model.GadgetReport{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&gadget, model.GadgetReport{RID: idInt}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &gadget
}
