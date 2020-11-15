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

// GetVote returns a vote
func GetVote(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vars := mux.Vars(r)
	id := vars["vid"]
	vote := getVoteOr404(db, id, w, r)
	if vote == nil {
		return
	}
	RespondJSON(w, http.StatusOK, vote)
}

// GetVotes returns all votes
func GetVotes(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	announcements := []model.Announcement{}
	db.Find(&announcements)
	RespondJSON(w, http.StatusOK, announcements)
}

// GetVotesFromCaller returns all votes by caller
func GetVotesFromCaller(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	votes := []model.Vote{}

	if err := db.Find(&votes, model.Vote{Caller: user}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, votes)
}

// GetVotesByType returns all votes by type
func GetVotesByType(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vType := vars["type"]

	votes := []model.Vote{}

	if err := db.Find(&votes, model.Vote{Type: vType}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, votes)
}

// GetActiveVotes returns all votes by type
func GetActiveVotes(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	votes := []model.Vote{}

	if err := db.Find(&votes, model.Vote{Complete: false}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, votes)
}

// GetPastVotes returns all votes by type
func GetPastVotes(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	votes := []model.Vote{}

	if err := db.Find(&votes, model.Vote{Complete: true}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, votes)
}

// GetSuccessfulVotes returns all votes by type
func GetSuccessfulVotes(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	votes := []model.Vote{}

	if err := db.Find(&votes, model.Vote{Passed: true, Complete: true}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, votes)
}

// GetFailedVotes returns all votes by type
func GetFailedVotes(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	votes := []model.Vote{}

	if err := db.Find(&votes, model.Vote{Passed: false, Complete: true}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, votes)
}

// CallVote calls a new vote
func CallVote(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	vote := model.Vote{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vote); err != nil {
		fmt.Println("error decoding vote:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&vote).Error; err != nil {
		fmt.Println("error saving vote:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, vote)
	fmt.Println("New Vote Called By " + vote.Caller + ": " + vote.Title + strconv.Itoa(int(vote.VID)))
}

func getVoteOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Vote {
	vote := model.Vote{}
	idInt, _ := strconv.ParseInt(id, 10, 32)
	if err := db.First(&vote, model.Vote{VID: idInt}).Error; err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &vote
}
