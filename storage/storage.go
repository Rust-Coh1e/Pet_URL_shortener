package storage

import (
	"sync"
	"sync/atomic"
	"time"
)

type URL struct {
	ID        string
	Original  string
	Clicks    int64
	CreatedAt time.Time
}

type Storage struct {
	urls map[string]*URL
	sync.RWMutex
}

func NewStorage() *Storage {
	url := make(map[string]*URL)

	return &Storage{
		urls: url,
	}
}

func (s *Storage) Save(url URL) {
	s.Lock()
	defer s.Unlock()
	s.urls[url.ID] = &url
}

func (s *Storage) Get(id string) (*URL, bool) {
	s.RLock()
	defer s.RUnlock()

	res, ok := s.urls[id]

	return res, ok
}

func (s *Storage) IncrementClick(id string) {
	s.RLock()             // защищаем чтение из map
	url, ok := s.urls[id] // берём указатель
	s.RUnlock()           // сразу отпускаем
	if ok {
		atomic.AddInt64(&url.Clicks, 1) // атомик без мьютекса
	}
}
