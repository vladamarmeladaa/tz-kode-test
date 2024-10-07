package get_notes

import (
	"encoding/json"
	"log/slog"
	"net/http"

	response "tz_kode/internal/lib/response"
	auth "tz_kode/internal/services/auth"
)

type Response struct {
	Notes []Note
	response.Response
}

type Note struct {
	Id    string `json:"id"`
	Title string `json:"title_note"`
	Text  string `json:"text_note"`
}

type NoteGetter interface {
	GetAllNotes(userId string) (any, error)
}

func New(log *slog.Logger, noteGetter NoteGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/handlers/get_notes"
		log := log.With(
			slog.String("op", op),
		)

		userId := r.Context().Value(auth.ContextKeyUser).(string)

		notes, err := noteGetter.GetAllNotes(userId)

		if err != nil {
			msg := "database communication error"
			log.Error(msg)
			w.WriteHeader(http.StatusInternalServerError)
			response.ResponseError(w, msg)
			return
		}

		// err = json.NewEncoder(w).Encode(notes)
		// if err != nil {
		// 	log.Error("encoding error")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	response.ResponseError(w, msg)
		// 	return
		// }
		log.Info("all notes are displayed")

		w.Header().Set("Content-Type", "application/json")
		responseOK(w, notes)
	}
}

func responseOK(w http.ResponseWriter, notes any) {
	err := json.NewEncoder(w).Encode(Response{Response: response.OK(), Notes: notes.([]Note)})
	if err != nil {
		// msg := "encoding error"
		// log.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		// response.ResponseError(w, msg)
		return
	}
}
