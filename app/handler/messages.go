package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	model "fgl-backend/app/model"

	"github.com/jinzhu/gorm"
)

// GetMessages returns all messages
func GetMessages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	messages := []model.Message{}
	db.Find(&messages)
	respondJSON(w, http.StatusOK, messages)
}

// ReceivedMessage handles an incoming message
func ReceivedMessage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	message := model.Message{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		fmt.Println("error decoding message:", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&message).Error; err != nil {
		fmt.Println("error saving message:", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, message)
}

// ReturnMessageBody returns the body of a Message object
func ReturnMessageBody(m string) string {
	message := model.Message{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageBody := message.Body
	return messageBody
}

// ReturnMessageAuthor returns the author of a Message object
func ReturnMessageAuthor(m string) string {
	message := model.Message{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageAuthor := message.Author
	return messageAuthor
}

// ReturnMessageTime returns the time of a Message object
func ReturnMessageTime(m string) time.Time {
	message := model.Message{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageTime := message.Time
	return messageTime
}
