// Package binding provides extended sources of data binding.
package binding

import (
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/gorilla/websocket"
)

// StringCloser is an extension of the String interface that allows resources to be freed
// using the standard `Close()` method.
type StringCloser interface {
	binding.String
	io.Closer
}

type webSocketString struct {
	binding.String
	conn *websocket.Conn
}

// NewWebSocketString returns a `String` binding to a web socket server specified as `url`.
// The resulting string will be set to the content of the latest message sent through the socket.
// You should also call `Close()` on the binding once you are done to free the connection.
func NewWebSocketString(url string) (StringCloser, error) {
	s := binding.NewString()
	sock, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			_, p, err := sock.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			err = s.Set(string(p))
			if err != nil {
				fyne.LogError("Failed to set string from web socket", err)
			}
		}
	}()
	ret := &webSocketString{String: s, conn: sock}
	return ret, nil
}

func (s *webSocketString) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}
