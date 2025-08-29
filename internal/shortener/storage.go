package shortener

import (
	"errors"
	"sync"
)

type URLRecord struct {
	Code string
	URL  string
}

type Storage interface {
	Save(code, longURL string) error
	Get(code string) (string, error)
	FindByURL(longURL string) (string, bool)
}

type MemoryStore struct {
	mu        sync.RWMutex
	codeToURL map[string]string
	urlToCode map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (s *MemoryStore) Save(code, longURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.codeToURL[code]; ok {
		return errors.New("code already exist")
	}
	s.codeToURL[code] = longURL
	s.urlToCode[longURL] = code
	return nil
}

func (s *MemoryStore) Get(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urlToCode[code]
	if !ok {
		return "", errors.New("not found")
	}
	return url, nil
}

func (s *MemoryStore) FindByURL(longURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, ok := s.urlToCode[longURL]
	return code, ok
}
