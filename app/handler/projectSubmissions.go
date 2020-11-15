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

// GetProjectSubmission returns a single project report
func GetProjectSubmission(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)
	id := vars["rid"]
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	RespondJSON(w, http.StatusOK, project)
}

// GetProjectSubmissions returns all project reports
func GetProjectSubmissions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	reports := []model.ProjectSubmission{}
	db.Find(&reports)
	RespondJSON(w, http.StatusOK, reports)
}

// GetProjectSubmissionsFromUser gets all announcements from a specific user
func GetProjectSubmissionsFromUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	anns := []model.Announcement{}

	if err := db.Find(&anns, model.Announcement{Author: user}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, anns)
}

// MakeProjectSubmission makes a new project report
func MakeProjectSubmission(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	project := model.ProjectSubmission{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		fmt.Println("error decoding report:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		fmt.Println("error saving report:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, project)
	fmt.Println("New Gadget Report By " + project.Author + ": " + project.Title)
}

// UpdateProjectSubmission updates an announcement
func UpdateProjectSubmission(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]
	gr := getProjectOr404(db, id, w, r)
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

// DeleteProjectSubmission deletes an announcement
func DeleteProjectSubmission(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)

	id := vars["rid"]

	gr := getProjectOr404(db, id, w, r)
	if gr == nil {
		return
	}
	if err := db.Delete(&gr).Error; err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

func getProjectOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.ProjectSubmission {
	project := model.ProjectSubmission{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&project, model.ProjectSubmission{PID: idInt}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &project
}
