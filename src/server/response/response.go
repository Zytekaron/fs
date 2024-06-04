package response

import (
	"encoding/json"
	"io"
	"net/http"
)

func WriteData(w http.ResponseWriter, statusCode int, data []byte) error {
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	return err
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func WriteStream(w http.ResponseWriter, statusCode int, r io.Reader) error {
	w.WriteHeader(statusCode)
	_, err := io.Copy(w, r)
	return err
}

func WriteError(w http.ResponseWriter, statusCode int) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(map[string]any{
		"error": http.StatusText(statusCode),
	})
}

func WriteErrorReason(w http.ResponseWriter, statusCode int, reason string) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(map[string]any{
		"error":  http.StatusText(statusCode),
		"reason": reason,
	})
}
