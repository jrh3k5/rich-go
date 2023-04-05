package client

import "fmt"

type ErrClosedConnection struct {
	cause error
}

func (e *ErrClosedConnection) Error() string {
	return fmt.Sprintf("the connection is closed: %v", e.cause)
}
