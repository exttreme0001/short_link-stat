package middleware

import (
	"context"
	"net/http"
	"restapi/configs"
	"restapi/pkg/jwt"
	"strings"
)

type key string // делается чтобы не затирать другие значения в программе

const (
	ContextEmailKey key = "ContentEmailKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}
func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer") {
			writeUnauthed(w)
			return
		}
		token := strings.TrimPrefix(authedHeader, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token) //приходящие параметры это валидность и data
		if !isValid {
			writeUnauthed(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx) // для передачи контекста необходимо пересоздать запроc
		next.ServeHTTP(w, req)    //все handlers теперь обогащены необходимым контекстом
	})
}
