package sign_up_user

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	response "tz_kode/internal/lib/response"
	auth "tz_kode/internal/services/auth"

	"github.com/go-playground/validator/v10"
)

type UserDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSignUpper interface {
	CreateUser(login, password string) error
	CheckUserIfExist(login string) (string, error)
}

type Validator interface {
	Struct(s interface{}) error
}

func New(log *slog.Logger, userSignUpper UserSignUpper, valid Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/handlers/sign_up_user"
		log := log.With(
			slog.String("op", op),
		)

		body := r.Body
		unmarshalUser, err := io.ReadAll(body)
		if err != nil {
			msg := "reading body error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		var user UserDTO
		err = json.Unmarshal(unmarshalUser, &user)
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

		userId, err := userSignUpper.CheckUserIfExist(user.Login)
		if userId != "" {
			msg := "user already exist"
			log.Error(msg, err)
			w.WriteHeader(http.StatusBadRequest)
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

		hashPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			msg := "hash password error"
			log.Error(msg)
			w.WriteHeader(http.StatusBadRequest)
			response.ResponseError(w, msg)
			return
		}
		user.Password = hashPassword

		err = userSignUpper.CreateUser(user.Login, user.Password)
		if err != nil {
			msg := "database communication error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		err = json.NewEncoder(w).Encode(response.OK())
		if err != nil {
			msg := "encoding error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		log.Info("user signed up")
		w.WriteHeader(http.StatusCreated)

	}
}
