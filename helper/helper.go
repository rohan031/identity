package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rohan031/identity/custom"
	"github.com/rohan031/identity/services"
)

type constraint interface {
	services.Identity
}

func DecodeJson[T constraint](w http.ResponseWriter, r *http.Request) (T, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload T
	err := decoder.Decode(&payload)

	if err != nil {
		var syntaxError *json.SyntaxError

		switch {
		case errors.As(err, &syntaxError):
			message := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return payload, &custom.MalformedRequest{
				Status:  http.StatusBadRequest,
				Message: message,
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			message := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return payload, &custom.MalformedRequest{
				Status:  http.StatusBadRequest,
				Message: message,
			}

		case errors.Is(err, io.EOF):
			message := "Request body must not be empty"
			return payload, &custom.MalformedRequest{
				Status:  http.StatusBadRequest,
				Message: message,
			}

		default:
			log.Printf("Error decoding request body: %v\n", err)
			return payload, err
		}
	}

	return payload, nil
}

func EncodeJson[T any](w http.ResponseWriter, status int, data T) {
	jsonRes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(jsonRes)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func errorResponse(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload services.JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	EncodeJson(w, statusCode, payload)
}

func HandleError(w http.ResponseWriter, err error) {
	var mr *custom.MalformedRequest

	if errors.As(err, &mr) {
		errorResponse(w, mr, mr.Status)
		return
	}

	err = errors.New(http.StatusText(http.StatusInternalServerError))
	errorResponse(w, err, http.StatusInternalServerError)

}
