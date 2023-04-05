package ipc

import "fmt"

type ErrClosedPipe struct {
	cause error
}

func (e *ErrClosedPipe) Error() string {
	return fmt.Sprintf("pipe is closing or is closed: %v", e.cause)
}
