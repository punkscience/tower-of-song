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
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Config struct {
	MusicFolders []string `json:"music_folders"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
}

var (
	db     *sql.DB
	config Config
	scanMutex  sync.Mutex // Mutex for controlling music folder scans
	tokenMutex sync.Mutex // Mutex for protecting tokenStore

	tokenStore = make(map[string]struct{}) // simple in-memory token store
)

func generateToken() string {
	return fmt.Sprintf("token-%d", time.Now().UnixNano())
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	header := r.Header.Get("Authorization")
	token := header
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	tokenMutex.Lock()
	_, ok := tokenStore[token]
	tokenMutex.Unlock()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func login(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if creds.Username == config.Username && creds.Password == config.Password {
		token := generateToken()
		tokenMutex.Lock()
		tokenStore[token] = struct{}{}
		tokenMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&config)
}

func ensureDataDir() error {
	return os.MkdirAll("/app/data", 0755)
}

func initDB() error {
	if err := ensureDataDir(); err != nil {
		return err
	}
	var err error
	db, err = sql.Open("sqlite", "/app/data/towerofsong.db")
	if err != nil {
		return err
	}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		fmt.Println("Error enabling WAL mode:", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS music (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE,
		title TEXT,
		artist TEXT,
		album TEXT,
		favourited INTEGER DEFAULT 0
	)`)
	if err != nil {
		return err
	}

	// Add 'favourited' column if it doesn't exist (for migrations)
	_, err = db.Exec("ALTER TABLE music ADD COLUMN favourited INTEGER DEFAULT 0")
	if err != nil {
		// Ignore error if column already exists
		if !strings.Contains(err.Error(), "duplicate column name") {
			return err
		}
	}

	return nil
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !requireAuth(w, r) {
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
	if !requireAuth(w, r) {
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
	if !requireAuth(w, r) {
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
	if !requireAuth(w, r) {
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
	log.Println("Scanning music folders...")
	scanMutex.Lock()
	defer scanMutex.Unlock()

	for _, folder := range config.MusicFolders {
		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(path, ".mp3") || strings.HasSuffix(path, ".flac") || strings.HasSuffix(path, ".wav")) {
				storeMetadata(path)
			}
			return nil
		})
	}

	// After scanning, clean up missing files from the database
	log.Println("Checking for missing files in database...")
	rows, err := db.Query("SELECT id, path FROM music")
	if err != nil {
		log.Println("Error querying music table for cleanup:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var path string
		if err := rows.Scan(&id, &path); err != nil {
			log.Println("Error scanning row during cleanup:", err)
			continue
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("File missing on disk, removing from database: %s (id=%d)\n", path, id)
			_, delErr := db.Exec("DELETE FROM music WHERE id = ?", id)
			if delErr != nil {
				log.Printf("Error deleting %s (id=%d) from database: %v\n", path, id, delErr)
			}
		}
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

func getTrackInfo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !requireAuth(w, r) {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT id, path, title, artist, album, favourited FROM music WHERE id = ?", id)
	var tid int
	var path, title, artist, album string
	var favourited bool
	err := row.Scan(&tid, &path, &title, &artist, &album, &favourited)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     fmt.Sprint(tid),
		"path":   path,
		"title":  title,
		"artist": artist,
		"album":  album,
		"favourited": favourited,
	})
}

func favouriteTrack(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !requireAuth(w, r) {
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("UPDATE music SET favourited = NOT favourited WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to update track", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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

	r := gin.Default()

	// Serve the test client as the homepage and as /index.html
	r.StaticFile("/", "templates/index.html")
	r.StaticFile("/index.html", "templates/index.html")

	// Serve static assets if needed (e.g., /templates/)
	r.Static("/templates", "templates")
	
	// Serve static CSS files
	r.Static("/static", "static")

	// Serve favicon
	r.StaticFile("/favicon.ico", "static/favicon.ico")

	// CORS middleware for API endpoints
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// API endpoints
	r.POST("/login", gin.WrapF(login))
	r.GET("/stats", gin.WrapF(getStats))
	r.GET("/list", gin.WrapF(listFiles))
	r.GET("/search", gin.WrapF(searchFiles))
	r.GET("/stream", gin.WrapF(streamFile))
	r.GET("/trackinfo", gin.WrapF(getTrackInfo))
	r.POST("/favourite", gin.WrapF(favouriteTrack))
	r.GET("/favourites", gin.WrapF(listFavourites))

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}

func listFavourites(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !requireAuth(w, r) {
		return
	}
	rows, _ := db.Query("SELECT id, path, title, artist, album FROM music WHERE favourited = 1 ORDER BY artist ASC, title ASC")
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
