package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed template.html
var embeddedTemplate string

type Link struct {
	URL         string `json:"url"`
	Image       string `json:"image"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type Config struct {
	Title              string   `json:"title"`
	Backgrounds        []string `json:"backgrounds"`
	Links              []Link   `json:"main_links"`
	AdditionalLinks    []Link   `json:"additional_links"`
	BackgroundColor    string   `json:"background_color"`
	BackgroundInterval int      `json:"background_interval"`
	FontColorPrimary   string   `json:"font_color1"`
	FontColorSecondary string   `json:"font_color2"`
}

var (
	cacheDir      = "/tmp/launcher-cache"
	allowedHashes = make(map[string]string)
	cacheMuxes    = sync.Map{}
)

func main() {
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "0.0.0.0:8080"
	}

	configFilePath := os.Getenv("CONFIG_FILEPATH")
	if configFilePath == "" {
		if _, err := os.Stat("config.json"); err == nil {
			configFilePath = "config.json"
		}
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatalf("Failed to create cache directory: %v", err)
	}

	config, err := loadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Critical error: failed to load configuration: %v", err)
	}

	tmpl, err := setupTemplate()
	if err != nil {
		log.Fatalf("Critical error: failed to parse HTML template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, config); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
		}
	})

	http.HandleFunc("/image", proxyImageHandler)

	server := &http.Server{
		Addr:         listenAddress,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server at %s", listenAddress)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func loadConfig(filepath string) (*Config, error) {
	defaultConfig := &Config{
		Title:              "Launcher",
		BackgroundInterval: 10,
		BackgroundColor:    "#6495ed",
		FontColorPrimary:   "#333",
		FontColorSecondary: "#555",
		Links: []Link{
			{
				URL:         "https://google.com",
				Image:       "https://placehold.co/150x150",
				Label:       "Google",
				Description: "Search engine",
			},
		},
	}

	if filepath == "" {
		return defaultConfig, nil
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(fileContent, defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	for _, link := range defaultConfig.Links {
		allowedHashes[hashURL(link.Image)] = link.Image
	}
	for _, link := range defaultConfig.AdditionalLinks {
		allowedHashes[hashURL(link.Image)] = link.Image
	}
	for _, bg := range defaultConfig.Backgrounds {
		allowedHashes[hashURL(bg)] = bg
	}

	return defaultConfig, nil
}

func setupTemplate() (*template.Template, error) {
	return template.New("template.html").Funcs(template.FuncMap{
		"hashURL":         hashURL,
		"replaceNewlines": replaceNewlines,
		"isURL":           isURL,
	}).Parse(embeddedTemplate)
}

func hashURL(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}

func replaceNewlines(input string) template.HTML {
	return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(input), "\n", "<br>"))
}

func isURL(image string) bool {
	return strings.HasPrefix(image, "http://") || strings.HasPrefix(image, "https://") || strings.HasPrefix(image, "file://")
}

func proxyImageHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "Bad Request: missing hash", http.StatusBadRequest)
		return
	}

	if _, ok := allowedHashes[hash]; !ok {
		http.Error(w, "Forbidden: hash not allowed", http.StatusForbidden)
		return
	}

	cachedPath, err := fetchAndCacheImage(hash)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error caching image: %v", err)
		return
	}

	file, err := os.Open(cachedPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error opening cached file: %v", err)
		return
	}
	defer file.Close()

	var length uint16
	if err := binary.Read(file, binary.BigEndian, &length); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error reading content type length: %v", err)
		return
	}

	contentType := make([]byte, length)
	if _, err := file.Read(contentType); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error reading content type: %v", err)
		return
	}

	w.Header().Set("Content-Type", string(contentType))
	if _, err := io.Copy(w, file); err != nil {
		log.Printf("Error serving file content: %v", err)
	}
}

func fetchAndCacheImage(hash string) (string, error) {
	source, ok := allowedHashes[hash]
	if !ok {
		return "", fmt.Errorf("hash not allowed: %s", hash)
	}

	cachedPath := filepath.Join(cacheDir, hash)

	if _, err := os.Stat(cachedPath); err == nil {
		return cachedPath, nil
	}

	muxIface, _ := cacheMuxes.LoadOrStore(hash, &sync.Mutex{})
	mux := muxIface.(*sync.Mutex)

	mux.Lock()
	defer mux.Unlock()

	var reader io.ReadCloser
	var contentType string

	if strings.HasPrefix(source, "file://") {
		localFilePath := strings.TrimPrefix(source, "file://")
		file, err := os.Open(localFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to open local file: %w", err)
		}
		reader = file
		defer reader.Close()

		buffer := make([]byte, 512)
		n, _ := reader.Read(buffer)
		if n > 0 {
			contentType = http.DetectContentType(buffer[:n])
			_, _ = file.Seek(0, io.SeekStart) // Reset file pointer to the start
		} else {
			contentType = "application/octet-stream"
		}
	} else {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Get(source)
		if err != nil {
			return "", fmt.Errorf("failed to fetch image: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
		}
		reader = resp.Body
		defer reader.Close()

		contentType = resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	file, err := os.Create(cachedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	length := uint16(len(contentType))
	if err := binary.Write(file, binary.BigEndian, length); err != nil {
		return "", fmt.Errorf("failed to write content type length: %w", err)
	}

	if _, err := file.Write([]byte(contentType)); err != nil {
		return "", fmt.Errorf("failed to write content type: %w", err)
	}

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to save content to cache: %w", err)
	}

	return cachedPath, nil
}
