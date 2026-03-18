package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

type BasicMessage struct {
	Message string `json:"message"`
}

type StandardMessage struct {
	Message interface{}
	Status  int `json:"status"`
}

type StandardError struct {
	Message     string   `json:"message"`
	Status      int      `json:"status"`
	Description []string `json:"description"`
}

func (e StandardError) Error() string {
	return fmt.Sprintf("Api error: %d", e.Status)
}

type Response interface {
	SendResponse(w http.ResponseWriter, r *http.Request)
	Test(http.Handler)
}

func (message *StandardMessage) SendReponse(
	w http.ResponseWriter,
	r *http.Request,
) {
	setHeader(&w, message.Status)

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/yaml":

		err := yaml.NewEncoder(w).Encode(message.Message)
		if err != nil {
			log.Error("Failed to Encode YAML", "error", err)
			return
		}
	default:

		err := json.NewEncoder(w).Encode(message.Message)
		if err != nil {
			log.Error("Failed to Encode")
			return
		}
	}
}

func (message *StandardError) SendReponse(
	w http.ResponseWriter,
	r *http.Request,
) {
	setHeader(&w, message.Status)

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/yaml":

		err := yaml.NewEncoder(w).Encode(message.Message)
		if err != nil {
			log.Error("Failed to Encode YAML", "error", err)
			return
		}
	default:

		err := json.NewEncoder(w).Encode(message.Message)
		if err != nil {
			log.Error("Failed to Encode")
			return
		}
	}
}

// NewFailureResponse returns a new instance of StandardError with the given
// message, status code and description.
func NewFailureResponse(
	message string,
	status int,
	description []string,
) StandardError {
	return StandardError{
		Message:     message,
		Status:      status,
		Description: description,
	}
}

// NewMessageResponse returns a new instance of StandardMessage with the given
// message and status code.
func NewMessageResponse(message interface{}, status int) *StandardMessage {
	return &StandardMessage{
		Message: message,
		Status:  status,
	}
}

// NewBasicResponse returns a new basic instance of StandardMessage with the
// message of OK and Status OK.
func NewBasicResponse() *StandardMessage {
	return &StandardMessage{
		Message: BasicMessage{
			Message: "OK",
		},
		Status: http.StatusOK,
	}
}

// SetHeader sets the HTTP response headers for a JSON response.
func setHeader(w *http.ResponseWriter, statusCode int) {
	(*w).WriteHeader(statusCode)
	(*w).Header().Set("Content-Type", "application/json")
}

type APIFunc func(
	w http.ResponseWriter,
	r *http.Request,
) (interface{}, error)

func Make(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := handler(w, r)
		if err != nil {
			var apiErr StandardError
			if errors.As(err, &apiErr) {
				apiErr.SendReponse(w, r)
				return
			}
			resp := NewFailureResponse(
				"Internal Server Error",
				http.StatusInternalServerError,
				[]string{err.Error()},
			)
			resp.SendReponse(w, r)
			return
		}

		message, ok := resp.(*StandardMessage)
		if ok {
			message.SendReponse(w, r)
			return
		}

		templMessage, ok := resp.(templ.Component)
		if ok {
			templMessage.Render(r.Context(), w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}
