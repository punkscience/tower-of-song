package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bogem/id3v2"
	_ "modernc.org/sqlite"
)

type Config struct {
	MusicFolders []string `json:"music_folders"`
}

var (
	db     *sql.DB
	config Config
	mutex  sync.Mutex
)

func loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&config)
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		return err
	}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		fmt.Println("Error enabling WAL mode:", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE music (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE,
		title TEXT,
		artist TEXT,
		album TEXT
	)`)
	return err
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM music").Scan(&count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"total_files": count})
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	rows, _ := db.Query("SELECT id, path, title, artist, album FROM music ORDER BY artist ASC, title ASC")
	var files []map[string]string

	for rows.Next() {
		var id int
		var path, title, artist, album string
		rows.Scan(&id, &path, &title, &artist, &album)
		files = append(files, map[string]string{"id": fmt.Sprint(id), "path": path, "title": title, "artist": artist, "album": album})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func searchFiles(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	query := r.URL.Query().Get("q")
	rows, _ := db.Query("SELECT id, path, title, artist, album FROM music WHERE title LIKE ? OR artist LIKE ? OR album LIKE ? OR path LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	var files []map[string]string
	for rows.Next() {
		var id int
		var path, title, artist, album string
		rows.Scan(&id, &path, &title, &artist, &album)
		files = append(files, map[string]string{"id": fmt.Sprint(id), "path": path, "title": title, "artist": artist, "album": album})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func getFilePathFromId(id string) string {
	var path string
	db.QueryRow("SELECT path FROM music WHERE id = ?", id).Scan(&path)
	return path
}

func streamFile(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	id := r.URL.Query().Get("id")

	filePath := getFilePathFromId(id)

	if filePath == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	log.Println("Streaming file: ", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(w, file)
}

func scanMusicFolders() {
	mutex.Lock()
	defer mutex.Unlock()

	for _, folder := range config.MusicFolders {
		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			log.Println("Scanning: ", folder)
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(path, ".mp3") || strings.HasSuffix(path, ".flac") || strings.HasSuffix(path, ".wav")) {
				storeMetadata(path)
			}
			return nil
		})
	}
}

func storeMetadata(path string) {
	title := filepath.Base(path)
	artist := "Unknown"
	album := "Unknown"

	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		log.Println("Error reading tags: ", err)
		log.Println("File was: ", path)
	} else {
		title = strings.TrimSpace(tag.Title())
		artist = strings.TrimSpace(tag.Artist())
		album = strings.TrimSpace(tag.Album())
	}
	defer tag.Close()

	if title == "" {
		title = filepath.Base(path)
	}
	if artist == "" {
		artist = "Unknown"
	}
	if album == "" {
		album = "Unknown"
	}

	_, err = db.Exec("INSERT OR IGNORE INTO music (path, title, artist, album) VALUES (?, ?, ?, ?)", path, title, artist, album)
	if err != nil {
		fmt.Println("DB insert error:", err)
	}
}

func main() {
	fmt.Println("Starting Tower of Song server...")

	if err := loadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	if err := initDB(); err != nil {
		fmt.Println("Error initializing DB:", err)
		return
	}

	go func() {
		for {
			scanMusicFolders()
			time.Sleep(24 * time.Hour)
		}
	}()

	http.HandleFunc("/stats", getStats)
	http.HandleFunc("/list", listFiles)
	http.HandleFunc("/search", searchFiles)
	http.HandleFunc("/stream", streamFile)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
