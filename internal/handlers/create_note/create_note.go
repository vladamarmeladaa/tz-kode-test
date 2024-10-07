package create_note

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

type Response struct {
	NoteId string `json:"noteId,omitempty"`
	response.Response
}

type NoteDTO struct {
	Title string `json:"title_note" validate:"required"`
	Text  string `json:"text_note" validate:"required"`
}

type SpellerResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=NoteCreator
type NoteCreator interface {
	CreateNote(title, text, userId string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SpellerValidator
type SpellerValidator interface {
	Validate(texts []string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Validator
type Validator interface {
	Struct(s interface{}) error
}

func New(log *slog.Logger, noteCreator NoteCreator, spellerValidator SpellerValidator, valid Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/handlers/create_note"
		log := log.With(
			slog.String("op", op),
		)

		userId := r.Context().Value(auth.ContextKeyUser).(string)

		body := r.Body
		b, err := io.ReadAll(body)
		if err != nil {
			msg := "reading body error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}
		defer r.Body.Close()

		var noteDTO NoteDTO
		err = json.Unmarshal(b, &noteDTO)
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
		err = valid.Struct(noteDTO)
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

		err = spellerValidator.Validate([]string{noteDTO.Title, noteDTO.Text})
		if err != nil {
			msg := fmt.Sprintf("speller validation error: %v", err)
			log.Error(msg, slog.String("inOp", "internal/services/speller"))
			w.WriteHeader(http.StatusBadRequest)
			response.ResponseError(w, msg)
			return
		}

		noteId, err := noteCreator.CreateNote(noteDTO.Title, noteDTO.Text, userId)
		if err != nil {
			msg := "database communication error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		// err = json.NewEncoder(w).Encode(noteId)
		// if err != nil {
		// 	log.Error("encoding error")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	response.Error(w, err.Error())
		// 	return
		// }

		log.Info("note added", slog.String("id", noteId))

		w.Header().Set("Content-Type", "application/json")
		responseOK(w, noteId)
	}
}

func responseOK(w http.ResponseWriter, noteId string) {
	err := json.NewEncoder(w).Encode(Response{Response: response.OK(), NoteId: noteId})
	if err != nil {
		// msg := "encoding error"
		// log.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		// response.ResponseError(w, msg)
		return
	}
}
