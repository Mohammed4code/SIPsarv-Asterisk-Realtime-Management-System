package main

import (
	"log"
	"net/http"
	"time"

	"sipsarv/internal/config"
	"sipsarv/internal/database"
	"sipsarv/internal/handlers"
	"sipsarv/internal/middleware"
)

func main() {
	// تحميل الإعدادات
	cfg := config.Load()

	// الاتصال بقاعدة البيانات
	if err := database.Connect(cfg.DB); err != nil {
		log.Printf("⚠️ تحذير: فشل الاتصال بقاعدة البيانات: %v", err)
		log.Printf("💡 تأكد من تشغيل MySQL على %s:%s", cfg.DB.Host, cfg.DB.Port)
	} else {
		log.Println("✅ تم الاتصال بقاعدة البيانات بنجاح")
	}
	defer database.Close()

	// إعداد المسارات
	mux := http.NewServeMux()

	// نقاط نهاية الصحة والفحص
	mux.HandleFunc("/api/health", middleware.CORS(handlers.Health))
	mux.HandleFunc("/api/asterisk/status", middleware.CORS(handlers.AsteriskStatus))
	mux.HandleFunc("/api/db/status", middleware.CORS(handlers.DBStatus))

	// نقاط نهاية الامتدادات
	mux.HandleFunc("/api/extensions/add", middleware.CORS(handlers.AddExtension))
	mux.HandleFunc("/api/extensions/list", middleware.CORS(handlers.ListExtensions))
	mux.HandleFunc("/api/extensions/delete", middleware.CORS(handlers.DeleteExtension))

	// نقاط نهاية جهات الاتصال
	mux.HandleFunc("/api/contacts/list", middleware.CORS(handlers.ListContacts))
	mux.HandleFunc("/api/contacts/delete", middleware.CORS(handlers.DeleteContact))

	// إعدادات Asterisk
	mux.HandleFunc("/api/asterisk/config", middleware.CORS(handlers.AsteriskConfigHandler))

	// الواجهة الأمامية
	mux.HandleFunc("/", handlers.ServeFrontend)

	// تشغيل الخادم
	port := cfg.Server.Port
	log.Printf("🚀 الخادم يعمل على http://localhost%s", port)
	log.Printf("📋 نقاط النهاية API متاحة على http://localhost%s/api/", port)
	log.Printf("⏰ بدأ التشغيل في: %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("❌ فشل تشغيل الخادم:", err)
	}
}