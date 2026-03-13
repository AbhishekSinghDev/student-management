package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AbhishekSinghDev/student-management/internal/types"
	"github.com/AbhishekSinghDev/student-management/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		decodeErr := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(decodeErr, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decodeErr))
			return
		}

		if decodeErr != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decodeErr))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"status": "OK"})
	}
}
