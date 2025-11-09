package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	// 🎵 Serve audio files
	r.HandleFunc("/audio/{filename}", serveAudio).Methods("GET")

	// ℹ️ Provide info about available audio files
	r.HandleFunc("/info", getAudioInfo).Methods("GET")

	// ✅ CORS setup (allow localhost + GitHub Pages)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://yourusername.github.io",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func serveAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := filepath.Join("audio", filename+".m4a")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "audio/mp4")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.m4a\"", filename))
	http.ServeFile(w, r, filePath)
}

// getAudioInfo lists all .m4a files in the audio directory
func getAudioInfo(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("audio")
	if err != nil {
		http.Error(w, "Could not read audio directory", http.StatusInternalServerError)
		return
	}

	type AudioFile struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	var audioList []AudioFile
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".m4a" {
			name := file.Name()
			audioList = append(audioList, AudioFile{
				Name: name,
				URL:  fmt.Sprintf("/audio/%s", name[:len(name)-4]), // strip ".m4a"
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audioList)
}
