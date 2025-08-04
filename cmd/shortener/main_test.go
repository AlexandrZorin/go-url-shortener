package main

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testURLs = map[string]string{
	"practicum": "https://practicum.yandex.ru",
	"google":    "https://google.com",
	"yandex":    "https://yandex.ru",
	"vk":        "https://vk.com",
}

func TestPostHandler(t *testing.T) {
	type want struct {
		status      int
		contentType string
		body        string
	}

	tests := []struct {
		name        string
		body        string
		contentType string
		method      string
		want        want
	}{
		{
			name:        "Valid case",
			body:        testURLs["practicum"],
			contentType: "text/plain",
			method:      http.MethodPost,
			want: want{
				status:      http.StatusCreated,
				contentType: "text/plain",
				body:        "http://localhost:8080/",
			},
		},
		{
			name:        "Wrong method",
			body:        testURLs["yandex"],
			contentType: "text/plain",
			method:      http.MethodGet,
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:        "Wrong Content-Type",
			body:        testURLs["google"],
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:        "Empty body",
			body:        "",
			contentType: "text/plain",
			method:      http.MethodPost,
			want: want{
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			postHandler(w, req)

			assert.Equal(t, tt.want.status, w.Code)

			if tt.want.body != "" {
				assert.True(t, strings.HasPrefix(w.Body.String(), tt.want.body))
			}
		})

		t.Run("Duplicate URL test", func(t *testing.T) {
			//first request
			req1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testURLs["google"]))
			req1.Header.Set("Content-Type", "text/plain")
			w1 := httptest.NewRecorder()
			postHandler(w1, req1)
			firstShortURL := w1.Body.String()

			//second request
			req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testURLs["google"]))
			req2.Header.Set("Content-Type", "text/plain")
			w2 := httptest.NewRecorder()
			postHandler(w2, req2)
			secondShortURL := w2.Body.String()

			assert.Equal(t, firstShortURL, secondShortURL)
		})
	}
}

func TestGetHandler(t *testing.T) {
	type want struct {
		status int
		URL    string
	}

	testData := map[string]string{
		"edasdkjf": testURLs["practicum"],
		"riudjvhk": testURLs["google"],
		"cjkhakfl": testURLs["yandex"],
		"odkjrhjk": testURLs["vk"],
	}

	urlStorage.urls = make(map[string]string)
	maps.Copy(urlStorage.urls, testData)

	tests := []struct {
		name     string
		method   string
		shortURL string
		want     want
	}{
		{
			name:     "Valid case practicum",
			method:   http.MethodGet,
			shortURL: "edasdkjf",
			want: want{
				status: http.StatusTemporaryRedirect,
				URL:    testURLs["practicum"],
			},
		},
		{
			name:     "Wrong method",
			method:   http.MethodPost,
			shortURL: "riudjvhk",
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name:     "Wrong short URL",
			method:   http.MethodGet,
			shortURL: "123qwer",
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/"+tt.shortURL, nil)
			w := httptest.NewRecorder()
			getHandler(w, req)

			assert.Equal(t, tt.want.status, w.Code)

			if tt.want.URL != "" {
				assert.Equal(t, tt.want.URL, w.Header().Get("Location"))
			}
		})
	}
}
