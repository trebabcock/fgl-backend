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
	messages := []model.Message{}
	db.Find(&messages)
	respondJSON(w, http.StatusOK, messages)
}

// ReceivedMessage handles an incoming message
func ReceivedMessage(db *gorm.DB, m string) {
	fmt.Println("received message", string(m))
	message := model.Message{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error:", err)
	}

	if err := db.Save(&message).Error; err != nil {
		fmt.Println("error:", err)
		return
	}
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
func ReturnMessageTime(m string) string {
	message := model.Message{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageTime := message.Time
	return messageTime
}
