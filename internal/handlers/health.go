package handlers

import (
	"fmt"
	"net/http"
	"time"

	"sipsarv/internal/asterisk"
	"sipsarv/internal/config"
	"sipsarv/internal/database"
)

func Health(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, "✅ الخادم يعمل", map[string]interface{}{
		"status":  "healthy",
		"uptime":  time.Now().Unix(),
		"service": "SIPsarv Asterisk Manager",
	})
}

func AsteriskStatus(w http.ResponseWriter, r *http.Request) {
	connected := asterisk.CheckConnection()
	cfg := config.Load()
	host, port, _, _ := cfg.Asterisk.Get()
	
	status := "غير متصل"
	if connected {
		status = "متصل"
	}
	
	sendSuccess(w, fmt.Sprintf("حالة Asterisk: %s", status), map[string]interface{}{
		"host":      host,
		"port":      port,
		"connected": connected,
		"status":    status,
	})
}

func DBStatus(w http.ResponseWriter, r *http.Request) {
	connected := database.IsConnected()
	cfg := config.Load()
	host, port, name, _, _ := cfg.DB.Get()
	
	status := "غير متصل"
	if connected {
		status = "متصل"
	}
	
	sendSuccess(w, fmt.Sprintf("حالة قاعدة البيانات: %s", status), map[string]interface{}{
		"host":      host,
		"port":      port,
		"name":      name,
		"connected": connected,
		"status":    status,
	})
}