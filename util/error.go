package util

import (
	"fmt"
)

type Message struct {
	Message string `json:"msg,string"`
}

func Str(message string) Message {
	return Message{Message: message}
}

func Err(err error) Message {
	return Message{Message: fmt.Sprintf("%s", err)}
}
