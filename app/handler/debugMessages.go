package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "fgl-backend/app/model"

	"github.com/jinzhu/gorm"
)

// GetDebugMessages returns all debug messages
func GetDebugMessages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	messages := []model.DebugMessage{}
	db.Find(&messages)
	respondJSON(w, http.StatusOK, messages)
}

// ReceivedDebugMessage handles an incoming debug message
func ReceivedDebugMessage(db *gorm.DB, m string) {
	fmt.Println("received message", string(m))
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error:", err)
	}

	if err := db.Save(&message).Error; err != nil {
		fmt.Println("error:", err)
		return
	}
}

// ReturnDebugMessageBody returns the body of a DebugMessage object
func ReturnDebugMessageBody(m string) string {
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageBody := message.Body
	return messageBody
}

// ReturnDebugMessageAuthor returns the author of a DebugMessage object
func ReturnDebugMessageAuthor(m string) string {
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageAuthor := message.Author
	return messageAuthor
}

// ReturnDebugMessageTime returns the time of a DebugMessage object
func ReturnDebugMessageTime(m string) string {
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageTime := message.Time
	return messageTime
}
