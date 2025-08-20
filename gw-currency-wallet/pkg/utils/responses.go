package utils

import (
	"encoding/json"
	"maps"
	"net/http"
)

func InternalErrorResponse(w http.ResponseWriter, headers http.Header) {
	WriteJSON(w, http.StatusInternalServerError, JSONEnveloper{"error": "internal server error"}, headers)
}

func BadRequestResponse(w http.ResponseWriter, msg JSONEnveloper, headers http.Header) {
	WriteJSON(w, http.StatusBadRequest, msg, headers)
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
