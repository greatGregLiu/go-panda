package client

import "fmt"

type Error struct {
	Code    int
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("panda: %d %s: %s", e.Code, e.Err, e.Message)
}
