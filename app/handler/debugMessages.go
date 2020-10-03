package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "fgl-backend/app/model"

	"github.com/jinzhu/gorm"
)

func GetDebugMessages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	messages := []model.DebugMessage{}
	db.Find(&messages)
	respondJSON(w, http.StatusOK, messages)
}

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

func ReturnDebugMessageBody(m string) string {
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageBody := message.Body
	return messageBody
}

func ReturnDebugMessageAuthor(m string) string {
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageAuthor := message.Author
	return messageAuthor
}

func ReturnDebugMessageTime(m string) string {
	message := model.DebugMessage{}
	if err := json.Unmarshal([]byte(m), &message); err != nil {
		fmt.Println("error", err)
	}
	messageTime := message.Time
	return messageTime
}
