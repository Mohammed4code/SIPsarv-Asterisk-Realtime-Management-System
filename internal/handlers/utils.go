package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"sipsarv/internal/database"
)

func sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(database.Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(database.Response{
		Success:   false,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func validateContact(c database.Contact) string {
	if c.FullName == "" {
		return "⚠️ الاسم الكامل مطلوب"
	}
	if c.PhoneNumber == "" {
		return "⚠️ رقم الهاتف مطلوب"
	}
	if c.Pin == "" {
		return "⚠️ PIN مطلوب"
	}
	if len(c.Pin) < 4 || len(c.Pin) > 8 {
		return "⚠️ PIN يجب أن يكون بين 4 و 8 أرقام"
	}
	if c.SipUsername == "" {
		return "⚠️ اسم مستخدم SIP مطلوب"
	}
	if c.SipSecret == "" {
		return "⚠️ SIP Secret مطلوب"
	}
	return ""
}