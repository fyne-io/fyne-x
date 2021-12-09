// Package binding provides extended sources of data binding.
package binding

import (
	"net/http"

	"fyne.io/fyne/v2/data/binding"

	"github.com/gorilla/websocket"
)

type webSocketString struct {
	binding.String
	conn *websocket.Conn
	prev error
}

// NewWebSocketString returns a `String` binding to a web socket server specified as `url`.
// The resulting string will be set to the content of the latest message sent through the socket.
// You should also call `Close()` on the binding once you are done to free the connection.
func NewWebSocketString(url string) (StringCloser, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}

	ret := &webSocketString{String: binding.NewString(), conn: conn}
	go ret.readMessages()
	return ret, nil
}

func (s *webSocketString) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *webSocketString) Get() (string, error) {
	if err := s.prev; err != nil {
		return "", err
	}

	return s.String.Get()
}

func (s *webSocketString) readMessages() {
	for {
		_, p, err := s.conn.ReadMessage()
		s.prev = err    // if no error we clear the state
		if err != nil { // permanent (could be connection closed)
			return
		}

		_ = s.Set(string(p)) // we control s, Set will not error
	}
}
