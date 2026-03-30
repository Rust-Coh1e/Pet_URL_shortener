package storage

import (
	"github.com/hashicorp/golang-lru/v2"
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
	urls *lru.Cache[string, *URL]
}

func NewStorage(size int) *Storage {
	c, err := lru.New[string, *URL](size)

	if err != nil {
		panic(err)
	}

	return &Storage{urls: c}
}

func (s *Storage) Save(url URL) {
	s.urls.Add(url.ID, &url)
}

func (s *Storage) Get(id string) (*URL, bool) {

	res, ok := s.urls.Get(id)

	// res, ok := s.urls[id]

	return res, ok
}

func (s *Storage) IncrementClick(id string) {

	url, ok := s.urls.Get(id) // берём указатель

	if ok {
		atomic.AddInt64(&url.Clicks, 1) // атомик без мьютекса
	}
}
