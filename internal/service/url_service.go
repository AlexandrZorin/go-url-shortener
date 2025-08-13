package service

import (
	"maps"
	"math/rand"
	"sync"
)

type URLService struct {
	mu   sync.Mutex
	urls map[string]string
}

func NewURLService() *URLService {
	return &URLService{
		urls: make(map[string]string),
	}
}

func (s *URLService) AddTestURLs(testData map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.urls, testData)
}

func (s *URLService) GenerateShortURL() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *URLService) CreateShortURL(originalURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var shortKey string
	for k, v := range s.urls {
		if v == originalURL {
			shortKey = k
			return shortKey
		}
	}

	shortKey = s.GenerateShortURL()
	s.urls[shortKey] = originalURL
	return shortKey
}

func (s *URLService) GetOriginalURL(shortKey string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, exists := s.urls[shortKey]
	return originalURL, exists
}
