package error

import (
	"strconv"
	"time"
)

type ChatError struct {
	ChatTime time.Time
}

func (e *ChatError) Error() string {
	return "At " + e.ChatTime.String() + "  Chat Error ,Check your Address or Password "
}

type ExecuteError struct {
	ExecuteTime time.Time
	Status      string
	StatusCode  int
}

func (e *ExecuteError) Error() string {
	return "At " + e.ExecuteTime.String() +
		"Executing API Error, Error code is " +
		strconv.Itoa(e.StatusCode) + "\n" +
		"According to " + e.Status
}
