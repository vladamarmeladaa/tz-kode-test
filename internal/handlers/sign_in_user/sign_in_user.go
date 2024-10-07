package sign_in_user

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	response "tz_kode/internal/lib/response"
	auth "tz_kode/internal/services/auth"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Token string
	response.Response
}

// type UserDTO struct {
// 	Login    string `json:"login"`
// 	Password string `json:"password"`
// }

type UserDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSignInner interface {
	CheckUserIfExist(login string) (string, error)
	GetUserPassword(login string) (string, error)
}

type Validator interface {
	Struct(s interface{}) error
}

func New(log *slog.Logger, userSignInner UserSignInner, valid Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/handlers/sign_in_user"
		log := log.With(
			slog.String("op", op),
		)

		body := r.Body
		b, err := io.ReadAll(body)
		if err != nil {
			msg := "reading body error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		var user UserDTO
		err = json.Unmarshal(b, &user)
		if err != nil {
			switch e := err.(type) {
			case *json.UnmarshalTypeError:
				msg := fmt.Sprintf("invalid type for field %v: Got: %s Want: %v", e.Field, e.Value, e.Type)
				log.Error(msg)
				w.WriteHeader(http.StatusBadRequest)
				response.ResponseError(w, msg)
				return
			default:
				log.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// if err != nil {
		// 	msg := "decoding error"
		// 	log.Error(msg)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	response.ResponseError(w, msg)
		// 	return
		// }

		userId, err := userSignInner.CheckUserIfExist(user.Login)
		if userId == "" {
			msg := "user unauthorized: CheckUserIfExist"
			log.Error(msg, err)
			w.WriteHeader(http.StatusUnauthorized)
			response.ResponseError(w, msg)
			return
		}

		err = valid.Struct(user)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			var errorString string
			for _, val := range errors {
				errorString += fmt.Sprintf("Field: %v ", val.Field())
				errorString += fmt.Sprintf("Tag: %v ", val.Tag())
			}
			log.Error("validation error:", err)
			w.WriteHeader(http.StatusBadRequest)
			response.ResponseError(w, fmt.Sprintf("validation error: %v", errorString))
			return
		}

		userPassword, err := userSignInner.GetUserPassword(user.Login)
		if err != nil {
			msg := err.Error()
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(user.Password))
		if err != nil {
			msg := "user unauthorized: invalid password"
			log.Error(msg)
			w.WriteHeader(http.StatusUnauthorized)
			response.ResponseError(w, msg)
			return
		}

		t, err := auth.CreateToken(user.Login)
		if err != nil {
			msg := "token creation error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		// err = json.NewEncoder(w).Encode(t)
		// if err != nil {
		// 	log.Error("encoding error")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	response.ResponseError(w, msg)
		// 	return
		// }
		log.Info("user signed in")

		responseOK(w, t)
	}
}

func responseOK(w http.ResponseWriter, t string) {
	err := json.NewEncoder(w).Encode(Response{Response: response.OK(), Token: t})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
