package handlers

import (
	"encoding/json"
	"net/http"
	"time"
    "fmt"
	"sipsarv/internal/asterisk"
	"sipsarv/internal/database"
)

func AddExtension(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var contact database.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		sendError(w, "طلب غير صحيح: "+err.Error(), http.StatusBadRequest)
		return
	}

	if msg := validateContact(contact); msg != "" {
		sendError(w, msg, http.StatusBadRequest)
		return
	}

	if !database.IsConnected() {
		sendError(w, "قاعدة البيانات غير متصلة — تحقق من إعدادات MySQL", http.StatusInternalServerError)
		return
	}

	if err := asterisk.InsertExtensionRealtime(contact); err != nil {
		sendError(w, "فشل إضافة الامتداد: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// إعادة تحميل PJSIP في الخلفية
	go asterisk.ReloadPJSIP()

	contact.CreatedAt = time.Now().Format(time.RFC3339)
	sendSuccess(w, fmt.Sprintf("✅ تم إضافة الامتداد %s للموظف %s بنجاح", contact.SipUsername, contact.FullName), contact)
}

func ListExtensions(w http.ResponseWriter, r *http.Request) {
	if !database.IsConnected() {
		sendError(w, "قاعدة البيانات غير متصلة", http.StatusInternalServerError)
		return
	}

	extensions, err := asterisk.ListExtensions()
	if err != nil {
		sendError(w, "فشل قراءة الامتدادات: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, fmt.Sprintf("📋 %d امتداد", len(extensions)), map[string]interface{}{
		"count":      len(extensions),
		"extensions": extensions,
	})
}

func DeleteExtension(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "طلب غير صحيح: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		sendError(w, "رقم الامتداد مطلوب", http.StatusBadRequest)
		return
	}

	if !database.IsConnected() {
		sendError(w, "قاعدة البيانات غير متصلة", http.StatusInternalServerError)
		return
	}

	if err := asterisk.DeleteExtensionRealtime(req.ID); err != nil {
		sendError(w, "فشل الحذف: "+err.Error(), http.StatusInternalServerError)
		return
	}

	go asterisk.ReloadPJSIP()
	sendSuccess(w, fmt.Sprintf("✅ تم حذف الامتداد %s", req.ID), nil)
}