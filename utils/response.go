package utils

import (
	"encoding/json"
	"net/http"
)

// Response common renderer response struct
type Response struct {
	StatusCode int         `json:"-"`
	Data       interface{} `json:"data,omitempty"`
	Message    interface{} `json:"message,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

// Render writes responses to header
func (r *Response) Render(w http.ResponseWriter) error {

	response, err := json.Marshal(r)

	if err != nil {
		logger.Errorln("Failed to marshal: ", err)
		return err
	}

	if r.StatusCode != 0 {
		w.WriteHeader(r.StatusCode)
	}
	if _, err = w.Write(response); err != nil {
		return err
	}
	return nil
}
