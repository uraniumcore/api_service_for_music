package main

import (
	"encoding/json"
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
	r.HandleFunc("/info/{filename}", infoHandler).Methods("GET")

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// infoHandler writes the name and extension of a file (if it exists) in JSON
func infoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// List of extensions we support in the audio folder
	exts := []string{".m4a", ".mp3"}

	for _, ext := range exts {
		fp := filepath.Join("audio", filename+ext)
		if _, err := os.Stat(fp); err == nil {
			// Found the file; return JSON with name and extension (without dot)
			resp := map[string]string{
				"name":      filename,
				"extension": ext[1:],
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	http.Error(w, "File not found", http.StatusNotFound)
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