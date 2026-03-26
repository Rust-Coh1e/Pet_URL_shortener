package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"url-shortener/generator"
	// "url-shortener/ratelimit"
	"url-shortener/repository"
	"url-shortener/storage"
)

// POST /shorten        — принять оригинальный URL, вернуть короткий
// GET  /:id            — редирект на оригинальный URL + инкремент кликов

type Handler struct {
	storage *storage.Storage
	repo    *repository.DB
}

func NewHandler(s *storage.Storage, r *repository.DB) *Handler {
	return &Handler{
		storage: s,
		repo:    r,
	}
}

type ShortenTaskReq struct {
	URL string `json:"url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "Application/json")

	var req ShortenTaskReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// req.URL
	shortURL := generator.GenerateURL()

	newURL := storage.URL{
		ID:        shortURL,
		Original:  req.URL,
		Clicks:    0,
		CreatedAt: time.Now(),
	}

	h.storage.Save(newURL)

	if err := h.repo.Save(newURL); err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortURL)
	// Save
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")

	parts := strings.Split(r.URL.Path, "/")

	id := parts[1]

	res, ok := h.storage.Get(id)

	if !ok {
		var err error
		res, err = h.repo.Get(id) // = без двоеточия — переиспользуем внешний res
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if res != nil {
			h.storage.Save(*res)
		}
	}

	if res == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.repo.IncrementClick(id)
	h.storage.IncrementClick(id)

	http.Redirect(w, r, res.Original, http.StatusMovedPermanently)
}
