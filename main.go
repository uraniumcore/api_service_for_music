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
			"http://localhost:8081",
			"https://uraniumcore.github.io/web-5.0/",
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
	// Try supported extensions in order of preference
	exts := []string{".m4a", ".mp3"}
	var foundPath string
	var foundExt string

	for _, ext := range exts {
		p := filepath.Join("audio", filename+ext)
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			foundExt = ext
			break
		}
	}

	if foundPath == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set Content-Type based on extension
	switch foundExt {
	case ".m4a":
		w.Header().Set("Content-Type", "audio/mp4")
	case ".mp3":
		w.Header().Set("Content-Type", "audio/mpeg")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s%s\"", filename, foundExt))
	http.ServeFile(w, r, foundPath)
}

// getAudioInfo lists all .m4a files in the audio directory
func getAudioInfo(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("audio")
	if err != nil {
		http.Error(w, "Could not read audio directory", http.StatusInternalServerError)
		return
	}

	type AudioFile struct {
		Name       string   `json:"name"`
		URL        string   `json:"url"`
		Extensions []string `json:"extensions"`
	}

	// Map base name (without ext) to AudioFile
	audioMap := map[string]*AudioFile{}
	supportedExts := map[string]bool{".m4a": true, ".mp3": true}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !supportedExts[ext] {
			continue
		}
		base := file.Name()[:len(file.Name())-len(ext)]
		entry, ok := audioMap[base]
		if !ok {
			entry = &AudioFile{
				Name:       base,
				URL:        fmt.Sprintf("/audio/%s", base),
				Extensions: []string{},
			}
			audioMap[base] = entry
		}
		entry.Extensions = append(entry.Extensions, ext[1:]) // without dot
	}

	var audioList []AudioFile
	for _, v := range audioMap {
		audioList = append(audioList, *v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audioList)
}
