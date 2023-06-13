package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/BogPin/real-time-chat/backend/api/utils"
)

type tokenBody struct {
	Token string `json:"token"`
}

type contextKey string

const TokenPayloadKey contextKey = "tokenPayload"

func GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return "", errors.New("absent or bad authorization token")
	}
	return authHeader[1], nil
}

func GetTokenFromQuery(r *http.Request) (string, error) {
	token := r.URL.Query().Get("jwt")
	if token == "" {
		return "", errors.New("absent or bad authorization token")
	}
	return token, nil
}

func GetAuthMiddleware(getToken func(r *http.Request) (string, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getToken(r)
			if err != nil {
				WriteError(w, utils.NewHttpError(err, http.StatusUnauthorized))
				return
			}

			body := tokenBody{token}
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(body)
			resp, err := http.Post("http://localhost:8081/auth/validate", "application/json", buf)
			if err != nil {
				WriteError(w, utils.NewHttpError(err, http.StatusUnauthorized))
				return
			}
			defer resp.Body.Close()

			statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
			if !statusOK {
				io.Copy(w, resp.Body)
				return
			}

			var payload TokenPayload
			err = json.NewDecoder(resp.Body).Decode(&payload)
			if err != nil {
				WriteError(w, utils.NewHttpError(err, http.StatusUnauthorized))
				return
			}

			ctx := context.WithValue(r.Context(), TokenPayloadKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
