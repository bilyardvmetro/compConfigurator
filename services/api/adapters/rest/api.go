package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"comp_config.com/services/api/core"
)

func NewLoginHandler(log *slog.Logger, loginer core.Loginer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Credentials struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		var creds Credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid JSON: ", http.StatusUnauthorized)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error("problem with closing wordsClient", "error", err)
			}
		}()

		if creds.Name == "" || creds.Password == "" {
			http.Error(w, "Name and password are required", http.StatusUnauthorized)
			return
		}

		token, err := loginer.Login(creds.Name, creds.Password)
		if err != nil {
			if err.Error() == core.ErrUnauthorized.Error() {
				http.Error(w, "wrong login or password", http.StatusUnauthorized)
				return
			}
			log.Error("Error with login", "error", err)
			http.Error(w, "Unknown error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(token)); err != nil {
			log.Error("error with response", "error", err)
		}
	}
}
