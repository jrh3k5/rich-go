package client

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jrh3k5/rich-go/ipc"
)

var logged bool

// Login sends a handshake in the socket and returns an error or nil
func Login(clientid string) error {
	if !logged {
		payload, err := json.Marshal(Handshake{"1", clientid})
		if err != nil {
			return err
		}

		err = ipc.OpenSocket()
		if err != nil {
			return err
		}

		if _, sendErr := ipc.Send(0, string(payload)); sendErr != nil {
			return fmt.Errorf("failed to send initial handshake: %w", sendErr)
		}
	} else {
		return errors.New("client is already logged in")
	}

	logged = true

	return nil
}

func Logout() error {
	logged = false

	return ipc.CloseSocket()
}

// SetActivity sets the activity.
// This can return ErrClosedConnection if the underlying connection has been closed.
// It is advised that the client be logged out, logged back in, and the message re-submitted in that situation.
func SetActivity(activity Activity) error {
	if !logged {
		return fmt.Errorf("client is not logged in")
	}

	payload, err := json.Marshal(Frame{
		"SET_ACTIVITY",
		Args{
			os.Getpid(),
			mapActivity(&activity),
		},
		getNonce(),
	})

	if err != nil {
		return fmt.Errorf("failed to marshal the SET_ACTIVITY payload to JSON: %w", err)
	}

	if _, sendErr := ipc.Send(1, string(payload)); sendErr != nil {
		cause := sendErr
		var pipeClosedErr *ipc.ErrClosedPipe
		if errors.As(sendErr, &pipeClosedErr) {
			cause = &ErrClosedConnection{
				cause: cause,
			}
		}
		return fmt.Errorf("failed to send SET_ACTIVITY payload: %w", cause)
	}
	return nil
}

func getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}
