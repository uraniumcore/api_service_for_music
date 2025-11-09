package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Gopher!")
	
	r := mux.NewRouter()
	r.HandleFunc("/audio/{filename}", serveAudio).Methods("GET")

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// serveAudio serves a local .m4a file
func serveAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := filepath.Join("audio", filename+".m4a")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers for audio file
	w.Header().Set("Content-Type", "audio/mp4") // correct MIME type for .m4a
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.m4a\"", filename))

	http.ServeFile(w, r, filePath)
}