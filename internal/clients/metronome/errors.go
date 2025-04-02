package metronome

import (
	"encoding/json"
	"io"
	"net/http"
)

type ClientError struct {
	Message string `json:"message"`
}

func (e *ClientError) Error() string {
	return e.Message
}

func ParseClientError(response io.Reader) *ClientError {
	var err ClientError
	if err := json.NewDecoder(response).Decode(&err); err != nil {
		return nil
	}
	return &err
}

func UnwrapClientError(res *http.Response) (*ClientError, bool) {
	c := ParseClientError(res.Body)
	return c, c != nil
}
