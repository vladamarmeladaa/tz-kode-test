package create_note

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	mocks "tz_kode/internal/handlers/create_note/mocks"
	mockLogger "tz_kode/internal/lib/logger"
	response "tz_kode/internal/lib/response"

	// response "tz_kode/internal/lib/response"
	auth "tz_kode/internal/services/auth"

	"github.com/stretchr/testify/require"
)

const userId = "vl"

func TestCreateNoteHandler(t *testing.T) {
	testCases := []struct {
		name               string
		title_note         string
		text_note          string
		respError          string
		mockError          error
		mockSpellerError   error
		mockValidatorError error
	}{
		{
			name:       "Success",
			title_note: "солнце",
			text_note:  "дождик капает по лужам",
		},
		{
			name:               "Empty title",
			title_note:         "",
			text_note:          "дождик капает по лужам",
			respError:          "Error",
			mockValidatorError: errors.New("unexpected error"),
		},
		{
			name:               "Empty text",
			title_note:         "солнце",
			text_note:          "",
			respError:          "Error",
			mockValidatorError: errors.New("unexpected error"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run("success", func(t *testing.T) {
			// t.Parallel()

			noteCreatorMock := mocks.NewNoteCreator(t)
			spellerValidatorMock := mocks.NewSpellerValidator(t)
			validatorMock := mocks.NewValidator(t)

			// if tc.respError != "" {
				validatorMock.On("Struct", NoteDTO{tc.title_note, tc.text_note}).
					Return(tc.mockValidatorError).
					Once()
				spellerValidatorMock.On("Validate", []string{tc.title_note, tc.text_note}).
					Return(tc.mockSpellerError).
					Once()
				noteCreatorMock.On("CreateNote", tc.title_note, tc.text_note, userId).
					Return(string("1"), tc.mockError).
					Once()

			// }

			handler := New(mockLogger.NewMockLogger(), noteCreatorMock, spellerValidatorMock, validatorMock)

			input := fmt.Sprintf(`{"title_note": "%s", "text_note": "%s"}`, tc.title_note, tc.text_note)

			req, err := http.NewRequest(http.MethodPost, "/notes", bytes.NewReader([]byte(input)))
			req = req.WithContext(context.WithValue(context.Background(), auth.ContextKeyUser, userId))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			// type r struct {
			// noteId string
			// resp   response.Response
			// }
			// var resp r

			var resp response.Response
			// fmt.Println(body)
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})

	}
}
