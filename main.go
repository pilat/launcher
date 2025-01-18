package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	Backgrounds        []string `json:"backgrounds"`
	Links              []Link   `json:"main_links"`
	AdditionalLinks    []Link   `json:"additional_links"`
	BackgroundColor    string   `json:"background_color"`
	BackgroundInterval int      `json:"background_interval"`
	FontColorPrimary   string   `json:"font_color1"`
	FontColorSecondary string   `json:"font_color2"`
}

func loadConfig(filepath string) (*Config, error) {
	defaultConfig := &Config{
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
			{
				URL:         "https://facebook.com",
				Image:       "https://placehold.co/150x150",
				Label:       "Facebook",
				Description: "Social media",
			},
			{
				URL:         "https://twitter.com",
				Image:       "https://placehold.co/150x150",
				Label:       "Twitter",
				Description: "Social media",
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

	return defaultConfig, nil
}

func setupTemplate() (*template.Template, error) {
	return template.New("template.html").Funcs(template.FuncMap{
		"jsonify": func(v interface{}) (template.JS, error) {
			bytes, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(bytes), nil
		},
	}).Parse(embeddedTemplate)
}

func main() {
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "0.0.0.0:8080"
	}

	configFilePath := os.Getenv("CONFIG_FILEPATH")
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
