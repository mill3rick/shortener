package handlers

import (
	"fmt"
	shortuuid "github.com/lithammer/shortuuid/v4"
	"io"
	"net/http"
	"net/url"
)

const (
	base = "http://localhost:8080/"
)

var bodyString string
var shortID string

// Shortener функция — обработчик HTTP-запроса
func ShortenerHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		// Проверяем тип запроса - валидируем только POST:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	contentType := req.Header.Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		//Проверяем Content-Type - он должен быть text/plain
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		//проверяем, всё ли нормально с телом запроса
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyString = string(body)
	//fmt.Println(bodyString)
	_, err = url.ParseRequestURI(bodyString)
	if err != nil {
		//проверка валидности http(s)-ссылки
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	//Создаем короткий uuid для создания нашей короткой ссылки
	shortID = shortuuid.New()[:8]

	shortURL := base + shortID
	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte(shortURL))

}

func GetRedirectHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		// Проверяем тип запроса - валидируем только POST:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// Проверка, есть ли для id из пути запроса короткая ссылка:
	if req.URL.Path[1:] != shortID {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(res, req, bodyString, http.StatusTemporaryRedirect)
}
