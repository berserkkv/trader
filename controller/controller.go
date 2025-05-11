package controller

import (
	"fmt"
	logger "github.com/berserkkv/trader/tools/log"
	"log/slog"
	"net/http"
)

var log *slog.Logger

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func Register() {
	log = logger.Get()

	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}
