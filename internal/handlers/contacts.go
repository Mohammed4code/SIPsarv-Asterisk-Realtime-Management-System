package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"sipsarv/internal/asterisk"
	"sipsarv/internal/database"
)

func ListContacts(w http.ResponseWriter, r *http.Request) {
	if !database.IsConnected() {
		sendError(w, "قاعدة البيانات غير متصلة", http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query(`
		SELECT full_name, passport, address, phone_number, pin, sip_username,
		       sip_secret, email, connection_type, reg_timeout,
		       asterisk_notes, created_at
		FROM contacts
		ORDER BY created_at DESC
	`)
	if err != nil {
		sendError(w, "فشل قراءة جهات الاتصال: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	contacts := []database.Contact{}
	for rows.Next() {
		var c database.Contact
		var createdAt time.Time
		var passport, address, email, notes sql.NullString
		
		if err := rows.Scan(&c.FullName, &passport, &address, &c.PhoneNumber, &c.Pin, &c.SipUsername,
			&c.SipSecret, &email, &c.ConnectionType, &c.RegTimeout, &notes, &createdAt); err != nil {
			continue
		}
		
		c.Passport = passport.String
		c.Address = address.String
		c.Email = email.String
		c.AsteriskNotes = notes.String
		c.CreatedAt = createdAt.Format(time.RFC3339)
		contacts = append(contacts, c)
	}

	sendSuccess(w, fmt.Sprintf("📋 %d جهة اتصال", len(contacts)), map[string]interface{}{
		"count":    len(contacts),
		"contacts": contacts,
	})
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "طلب غير صحيح: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !database.IsConnected() {
		sendError(w, "قاعدة البيانات غير متصلة", http.StatusInternalServerError)
		return
	}

	// الحصول على sip_username من رقم الهاتف
	var sipUsername string
	err := database.DB.QueryRow("SELECT sip_username FROM contacts WHERE phone_number = ?", req.PhoneNumber).Scan(&sipUsername)
	if err != nil {
		sendError(w, "جهة الاتصال غير موجودة", http.StatusNotFound)
		return
	}

	if err := asterisk.DeleteExtensionRealtime(sipUsername); err != nil {
		sendError(w, "فشل الحذف: "+err.Error(), http.StatusInternalServerError)
		return
	}

	go asterisk.ReloadPJSIP()
	sendSuccess(w, "✅ تم حذف جهة الاتصال والامتداد المرتبط بها", nil)
}