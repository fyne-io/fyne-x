// Package binding provides extended sources of data binding.
package binding

import (
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/gorilla/websocket"
)

type webSocketString struct {
	binding.String
	conn *websocket.Conn
}

// NewWebSocketString returns a `String` binding to a web socket server specified as `url`.
// The resulting string will be set to the content of the latest message sent through the socket.
// You should also call `Close()` on the binding once you are done to free the connection.
func NewWebSocketString(url string) (StringCloser, error) {
	s := binding.NewString()
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			_, p, err := conn.ReadMessage()
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
	ret := &webSocketString{String: s, conn: conn}
	return ret, nil
}

func (s *webSocketString) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}
