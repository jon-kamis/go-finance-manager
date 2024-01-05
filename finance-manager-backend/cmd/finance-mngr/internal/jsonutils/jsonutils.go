package jsonutils

import (
	"encoding/json"
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"io"
	"net/http"
)

type JSONUtils interface {
	WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error
	ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error
	ErrorJSON(w http.ResponseWriter, err error, status ...int) error
}

type JSONUtil struct{}

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (j *JSONUtil) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if status != 200 {
		w.WriteHeader(status)
	}

	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func (j *JSONUtil) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	//dec.DisallowUnknownFields()

	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (j *JSONUtil) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true

	if err != nil {
		payload.Message = err.Error()
	} else {
		payload.Message = constants.JSONDefaultErrorMessage
	}

	return j.WriteJSON(w, statusCode, payload)
}
