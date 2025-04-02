package metronome

import (
	"encoding/json"
	"io"
)

type ClientError struct {
	Message string `json:"message"`
}

func IsClientError(response io.Reader, expectedMessage string) bool {
	var err ClientError
	if err := json.NewDecoder(response).Decode(&err); err != nil {
		return false
	}

	return err.Message == expectedMessage
}
