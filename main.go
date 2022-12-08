package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Redirections map[string]string `json:"redirections"`
	Bind         string            `json:"bind"`
}

var config Config
var configFile string

func loadConfig() error {
	configStr, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configStr, &config)
	if err != nil {
		return err
	}
	return nil
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Remove trailing slash in r.URL.Path
	if len(r.URL.Path) > 1 && r.URL.Path[len(r.URL.Path)-1] == '/' {
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
	}
	if url, ok := config.Redirections[r.URL.Path]; ok {
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	flag.StringVar(&configFile, "c", "/opt/goto/goto.json", "config file")
	flag.Parse()
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(config.Bind, nil))
}
