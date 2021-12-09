// Package binding provides extended sources of data binding.
package binding

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttString struct {
	binding.String
	conn  mqtt.Client
	topic string
}

// NewMqttString returns a `String` binding to a web socket server specified as `url`.
// The resulting string will be set to the content of the latest message sent through the socket.
// You should also call `Close()` on the binding once you are done to free the connection.
func NewMqttString(conn mqtt.Client, topic string) (StringCloser, error) {
	s := binding.NewString()

	go func() {
		token := conn.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
			s.Set(string(m.Payload()))
		})

		token.Wait()

		if token.Error() != nil {
			fyne.LogError("Failed to subscribe", token.Error())
		}
	}()

	ret := &mqttString{String: s, conn: conn, topic: topic}

	return ret, nil
}

func (s *mqttString) Set(val string) error {
	token := s.conn.Publish(s.topic, 0, false, val)

	go func() {
		token.Wait()

		if token.Error() != nil {
			fyne.LogError("Failed to publish", token.Error())
		}
	}()

	return nil
}

func (s *mqttString) Close() error {
	if s.conn == nil {
		return nil
	}

	s.conn.Unsubscribe(s.topic)
	s.conn = nil

	return nil
}
