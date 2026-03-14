package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/AbhishekSinghDev/student-management/internal/storage"
	"github.com/AbhishekSinghDev/student-management/internal/types"
	"github.com/AbhishekSinghDev/student-management/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		// req body validation
		validationError := validator.New().Struct(student)
		if validationError != nil {
			vError := validationError.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(vError))
			return
		}

		lastId, creationError := storage.CreateStudent(student.Name, student.Email, student.Age)
		if creationError != nil {
			response.WriteJson(w, http.StatusInternalServerError, creationError)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		parsedId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}

		student, err := storage.GetStudentById(int64(parsedId))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
