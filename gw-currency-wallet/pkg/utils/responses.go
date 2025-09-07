package utils

import (
	"encoding/json"
	"maps"
	"net/http"
)

func InternalErrorResponse(w http.ResponseWriter) {
	ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
}

func BadRequestResponse(w http.ResponseWriter, msg any) {
	ErrorResponse(w, http.StatusBadRequest, msg)
}

func UnprocessableEntityResponse(w http.ResponseWriter, msg any) {
	ErrorResponse(w, http.StatusUnprocessableEntity, msg)
}
func UnauthorizedResponse(w http.ResponseWriter, msg any) {
	ErrorResponse(w, http.StatusUnauthorized, msg)
}

// записывает данные в формат JSON и отправляет клиенту
func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// записывает хедеры
	maps.Copy(w.Header(), headers)
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ErrorResponse(w http.ResponseWriter, status int, message any) {
	enc := JSONEnveloper{"error": message}

	err := WriteJSON(w, status, enc, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
