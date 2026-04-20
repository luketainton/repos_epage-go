package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

const (
	errExecutingTemplate = "Error executing template: %v"
)

var (
	templateCache = make(map[string]*template.Template)
	cacheMutex    sync.RWMutex
)

func loadTemplate(baseDir, templateName string) (*template.Template, error) {
	cacheMutex.RLock()
	if tmpl, exists := templateCache[templateName]; exists {
		cacheMutex.RUnlock()
		return tmpl, nil
	}
	cacheMutex.RUnlock()

	templatePath := filepath.Join(baseDir, "templates", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	templateCache[templateName] = tmpl
	cacheMutex.Unlock()

	return tmpl, nil
}

func handleIndex(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexTemplate, err := loadTemplate(baseDir, "index.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := map[string]string{
			"status": "",
		}

		err = indexTemplate.Execute(w, data)
		if err != nil {
			log.Printf(errExecutingTemplate, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func handleSend(baseDir, apiToken, userKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexTemplate, err := loadTemplate(baseDir, "index.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Parse form data
		err = r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

	// Validate required fields
	if name == "" || email == "" || message == "" {
		data := map[string]string{
			"status": "fail",
		}
		err = indexTemplate.Execute(w, data)
		if err != nil {
			log.Printf(errExecutingTemplate, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

		// Send page via Pushover
		success, err := sendPage(apiToken, userKey, name, email, message)

		status := "fail"
		if success {
			status = "success"
		}

		data := map[string]string{
			"status": status,
		}

		if err != nil {
			log.Printf("Error sending page: %v", err)
		}

		err = indexTemplate.Execute(w, data)
		if err != nil {
			log.Printf(errExecutingTemplate, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
