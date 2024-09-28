package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSendResponseStandardMessage(t *testing.T) {
	message := &BasicMessage{
		Message: "Hello World!",
	}
	resp := &StandardMessage{
		Status:  http.StatusAccepted,
		Message: message,
	}
	w := httptest.NewRecorder()
	resp.SendReponse(w)

	// Check the status code is set correctly
	if w.Code != 202 {
		t.Errorf("Expected status code to be 202, but got %d", w.Code)
	}

	// Check that the content type header is set to "application/json"
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"Expected content type to be 'application/json', but got %s",
			contentType,
		)
	}

	// Check that the message is written to the response body correctly
	response := &BasicMessage{}
	err := json.NewDecoder(w.Body).Decode(response)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(message, response) {
		t.Errorf("Expected body to be %s but got %s", message, response)
	}
}

func TestSendResponseStandardError(t *testing.T) {
	resp := &StandardError{
		Status:      http.StatusInternalServerError,
		Message:     "An error has occured",
		Description: []string{"An Error"},
	}

	w := httptest.NewRecorder()
	resp.SendReponse(w)

	// Check the status code is set correctly
	if w.Code != 500 {
		t.Errorf("Expected status code to be 500, but got %d", w.Code)
	}

	// Check that the content type header is set to "application/json"
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf(
			"Expected content type to be 'application/json',"+
				" but got %s",
			contentType,
		)
	}

	// Check that the message is written to the response body correctly
	response := &StandardError{}
	err := json.NewDecoder(w.Body).Decode(response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Message != resp.Message {
		t.Errorf(
			"Expected Message of %s but got %s",
			resp.Message,
			response.Message,
		)
	}

	if !reflect.DeepEqual(resp, response) {
		t.Errorf("Expected Message of %v but got %v", resp, response)
	}
}

// NewFailureResponse returns a new instance of StandardError with the given
// message, status code and description.
func TestNewFailureResponse(t *testing.T) {
	message := "An error has occured"
	status := http.StatusInternalServerError
	description := []string{"An Error"}
	resp := NewFailureResponse(message, status, description)

	if resp.Status != status {
		t.Errorf(
			"Expected Status to be %d but got %d",
			status,
			resp.Status,
		)
	}
	if resp.Message != message {
		t.Errorf(
			"Expected Status to be %s but got %s",
			message,
			resp.Message,
		)
	}
	if !reflect.DeepEqual(description, resp.Description) {
		t.Errorf(
			"Expected Status to be %v but got %v",
			description,
			resp.Description,
		)
	}
}

func TestNewMessageResponse(t *testing.T) {

	message := &BasicMessage{
		Message: "Hello World!",
	}

	resp := NewMessageResponse(message, http.StatusOK)

	if resp.Status != http.StatusOK {
		t.Errorf(
			"Expected Status to be %d but got %d",
			http.StatusOK,
			resp.Status,
		)
	}
	if !reflect.DeepEqual(message, resp.Message) {
		t.Errorf(
			"Expected Message to be %s but got %s",
			message,
			resp.Message,
		)
	}
}
