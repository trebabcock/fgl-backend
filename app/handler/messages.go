package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "fgl-backend/app/model"

	"github.com/jinzhu/gorm"
)

// GetMessages returns all messages
func GetMessages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !Authorize(db, w, r) {
		return
	}
	messages := []model.Message{}
	db.Find(&messages)
	RespondJSON(w, http.StatusOK, messages)
}

// ReceivedMessage handles an incoming message
func ReceivedMessage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	message := model.Message{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		fmt.Println("error decoding message:", err)
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&message).Error; err != nil {
		fmt.Println("error saving message:", err)
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, message)
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
