package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var shortID2 string

func testRequest1(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, contentType string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", contentType)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)
	shortID2 = "/"
	shortID2 += string(respBody)[strings.LastIndex(string(respBody), "/")+1:]
	return resp
}

func testRequest2(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, io.NopCloser(strings.NewReader(" ")))
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp
}

func TestShortenerHandler(t *testing.T) {
	ts := httptest.NewServer(ShortenerRouter())
	defer ts.Close()
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
			testNumber:  "test 1: not allowed method",
			method:      http.MethodGet,
			contentType: "text/plain; charset=utf-8",
			body:        "https://google.com",
			want:        want{expectedCode: http.StatusMethodNotAllowed, expectedContentType: ""},
		},
		{
			testNumber:  "test 2: bad body",
			method:      http.MethodPost,
			contentType: "text/plain; charset=utf-8",
			body:        "a",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 3: empty body",
			method:      http.MethodPost,
			contentType: "text/plain; charset=utf-8",
			body:        "",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 4: bad content-type",
			method:      http.MethodPost,
			contentType: "application/x-www-form-urlencoded",
			body:        "https://google.com",
			want:        want{expectedCode: http.StatusBadRequest, expectedContentType: ""},
		},
		{
			testNumber:  "test 5: passed",
			method:      http.MethodPost,
			contentType: "text/plain; charset=utf-8",
			body:        "https://google.com",
			want:        want{expectedCode: http.StatusCreated, expectedContentType: "text/plain; charset=utf-8"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testNumber, func(t *testing.T) {
			resp := testRequest1(t, ts, tc.method, "/", io.NopCloser(strings.NewReader(tc.body)), tc.contentType)
			assert.Equal(t, tc.want.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tc.want.expectedContentType, resp.Header.Get("Content-Type"), "Тело не совпадает с ожидаемым, тело из результата: %s", resp.Header.Get("Content-Type"))
		})
	}
}

func TestGetRedirectHandler(t *testing.T) {
	ts := httptest.NewServer(ShortenerRouter())
	defer ts.Close()
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
				expectedCode:        http.StatusMethodNotAllowed,
				expectedContentType: "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testNumber, func(t *testing.T) {
			resp := testRequest2(t, ts, tc.method, tc.path)
			assert.Equal(t, tc.want.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			assert.Equal(t, tc.want.expectedContentType, resp.Header.Get("Content-Type"), "Хэдер не совпадает с ожидаемым, хэдер из результата: %s", resp.Header.Get("Content-Type"))
		})
	}
}
