package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dant1one/URL/storage"
	"github.com/gorilla/mux"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode(req.URL)

	if err := storage.SaveURL(shortCode, req.URL); err != nil {
		log.Printf("Error saving to database: %v", err)
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{
		ShortURL: fmt.Sprintf("http://localhost:8081/%s", shortCode),
	})
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	longURL, err := storage.GetURL(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func generateShortCode(longURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}

func main() {
	storage.InitDB()
	r := mux.NewRouter()
	r.HandleFunc("/shorten", ShortenURL).Methods("POST")
	r.HandleFunc("/{shortCode}", HandleRedirect).Methods("GET")

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
