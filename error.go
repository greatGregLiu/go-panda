package panda

import (
	"errors"
	"fmt"
)

var errClientNil = errors.New("the client cannot be nil")

type PandaError struct {
	Code    int
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e PandaError) Error() string {
	return fmt.Sprintf("panda: %d %s: %s", e.Code, e.Err, e.Message)
}

func clientError(cl Client) (err error) {
	if cl == nil {
		err = errClientNil
	}
	return
}
