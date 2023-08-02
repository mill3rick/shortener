// пакеты исполняемых приложений должны называться main
package main

import (
	"github.com/gorilla/mux"
	"github.com/mill3rick/shortener/internal/handlers"
	"net/http"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	//mux := http.NewServeMux()
	rtr := mux.NewRouter()
	rtr.HandleFunc(`/`, handlers.Shortener)
	rtr.HandleFunc(`/{id}`, handlers.GetRedirect)
	return http.ListenAndServe(`:8080`, rtr)
}
