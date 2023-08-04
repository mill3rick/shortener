package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var shortID2 string

func TestShortenerHandler(t *testing.T) {
	type want struct {
		expectedCode        int
		expectedContentType string
	}
	testCases := []struct {
		testNumber  string
		method      string
		contentType string
		body        string
		want        want
	}{
		{
			testNumber:  "test 1",
			method:      http.MethodGet,
			contentType: "text/plain; charset=utf-8",
			body:        "https://google.com",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 2",
			method:      http.MethodPost,
			contentType: "text/plain; charset=utf-8",
			body:        "a",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 3",
			method:      http.MethodPost,
			contentType: "text/plain; charset=utf-8",
			body:        "",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 4",
			method:      http.MethodPost,
			contentType: "application/x-www-form-urlencoded",
			body:        "https://google.com",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 5",
			method:      http.MethodPost,
			contentType: "text/plain; charset=utf-8",
			body:        "https://google.com",
			want:        want{expectedCode: http.StatusCreated, expectedContentType: "text/plain; charset=utf-8"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testNumber, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", io.NopCloser(strings.NewReader(tc.body)))
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()
			ShortenerHandler(w, r)
			body, _ := io.ReadAll(w.Body)
			shortID2 = "/"
			shortID2 += string(body)[strings.LastIndex(string(body), "/")+1:]
			assert.Equal(t, tc.want.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tc.want.expectedContentType, w.Header().Get("Content-Type"), "Тело не совпадает с ожидаемым, тело из результата: %s", w.Header().Get("Content-Type"))
		})
	}
}

func TestGetRedirectHandler(t *testing.T) {
	type want struct {
		expectedCode        int
		expectedContentType string
	}
	testCases := []struct {
		testNumber string
		method     string
		path       string
		want       want
	}{
		{
			testNumber: "test 1",
			method:     http.MethodGet,
			path:       shortID2,
			want: want{
				expectedCode:        http.StatusTemporaryRedirect,
				expectedContentType: "text/html; charset=utf-8",
			},
		},
		{
			testNumber: "test 2",
			method:     http.MethodGet,
			path:       "/nottruevalue",
			want: want{
				expectedCode:        http.StatusBadRequest,
				expectedContentType: "",
			},
		},
		{
			testNumber: "test 3",
			method:     http.MethodPost,
			path:       shortID2,
			want: want{
				expectedCode:        http.StatusBadRequest,
				expectedContentType: "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testNumber, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			// вызовем хендлер как обычную функцию, без запуска самого сервера
			GetRedirectHandler(w, r)

			assert.Equal(t, tc.want.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			assert.Equal(t, tc.want.expectedContentType, w.Header().Get("Content-Type"), "Тело не совпадает с ожидаемым, тело из результата: %s", w.Header().Get("Content-Type"))
		})
	}
}
