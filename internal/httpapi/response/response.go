package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func Fail(w http.ResponseWriter, statusCode int, msg string) {
	JSON(w, statusCode, Error{Error: msg})
}
