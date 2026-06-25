package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ServeFrontend(w http.ResponseWriter, r *http.Request) {
	// تجاهل طلبات API
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// تحديد المسار
	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "web/index.html"
	}

	// إذا كان الملف غير موجود، ارجع index.html (لـ SPA)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "web/index.html"
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// تحديد نوع المحتوى
	ext := filepath.Ext(filePath)
	contentType := "text/html"
	switch ext {
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".json":
		contentType = "application/json"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".svg":
		contentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write(content)
}