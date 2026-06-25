package handlers

import (
	"encoding/json"
	"net/http"

	"sipsarv/internal/config"
)

func AsteriskConfigHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()

	switch r.Method {
	case "GET":
		_, port, user, _ := cfg.Asterisk.Get()
		sendSuccess(w, "✅ إعدادات Asterisk الحالية", map[string]interface{}{
			"host":     cfg.Asterisk.Host,
			"port":     port,
			"user":     user,
			"password": "********",
		})

	case "POST":
		var req struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "طلب غير صحيح: "+err.Error(), http.StatusBadRequest)
			return
		}

		cfg.Asterisk.Set(req.Host, req.Port, req.User, req.Password)
		sendSuccess(w, "✅ تم تحديث إعدادات Asterisk", map[string]string{
			"host": cfg.Asterisk.Host,
			"port": cfg.Asterisk.Port,
			"user": cfg.Asterisk.User,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}