package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	response "tz_kode/internal/lib/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("secret")

const (
	ContextKeyUser string = "secretKeyUser"
)

type UserGetter interface {
	CheckUserIfExist(login string) (string, error)
}

func AuthMiddleware(next http.HandlerFunc, log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/services/authMiddleware"
		log := log.With(
			slog.String("op", op),
		)

		token := r.Header.Get("Authorization")
		if len(strings.Fields(token)) == 1 {
			// }
			// if token == "" {
			msg := "token is empty"
			log.Error(msg)
			w.WriteHeader(http.StatusUnauthorized)
			response.ResponseError(w, msg)
			return
		}

		payload, _ := jwt.Parse(strings.Fields(token)[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		})

		// payload := &jwt.Token{Raw: strings.Fields(token)[1]}

		login, err := payload.Claims.GetSubject()
		if err != nil {
			msg := "token parsing error"
			log.Error(msg)
			w.WriteHeader(http.StatusUnauthorized)
			response.ResponseError(w, msg)
			return
		}

		userId, err := userGetter.CheckUserIfExist(login)
		if userId == "" || err != nil {
			msg := "user not found"
			log.Error(msg)
			w.WriteHeader(http.StatusUnauthorized)
			response.ResponseError(w, msg)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUser, userId)

		r = r.WithContext(ctx)

		next(w, r)

	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CreateToken(login string) (string, error) {
	payload := jwt.MapClaims{
		"sub": login,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return t, nil

}
