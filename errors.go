package lambique

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents an error based on RFC7807.
type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// NewErrorResponse creates a new ErrorResponse.
func NewErrorResponse(r *http.Request, err error, status int) *ErrorResponse {
	return &ErrorResponse{
		Type:     "about:blank",
		Title:    fmt.Sprintf("%s", err),
		Status:   status,
		Detail:   fmt.Sprintf("%s", err),
		Instance: r.RequestURI,
	}
}

// JSON writes the ErrorResponse as JSON.
func (resp *ErrorResponse) JSON(w http.ResponseWriter) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(resp.Status)
	return json.NewEncoder(w).Encode(resp)
}

// MustJSON is like JSON but panics if the JSON encoder returns error.
func (resp *ErrorResponse) MustJSON(w http.ResponseWriter) {
	err := resp.JSON(w)
	if err != nil {
		panic(err)
	}
}
