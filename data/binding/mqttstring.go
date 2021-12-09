// Package binding provides extended sources of data binding.
package binding

import (
	"fyne.io/fyne/v2/data/binding"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttString struct {
	binding.String
	conn  mqtt.Client
	topic string
	err   error
}

// NewMqttString returns a `String` binding to a web socket server specified as `url`.
// The resulting string will be set to the content of the latest message sent through the socket.
// You should also call `Close()` on the binding once you are done to free the connection.
func NewMqttString(conn mqtt.Client, topic string) (StringCloser, error) {
	ret := &mqttString{String: binding.NewString(), conn: conn, topic: topic}

	token := conn.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		ret.String.Set(string(m.Payload()))
	})

	token.Wait()

	if err := token.Error(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *mqttString) Set(val string) error {
	token := s.conn.Publish(s.topic, 0, false, val)

	token.Wait()
	if s.err = token.Error(); s.err != nil {
		return s.err
	}
	s.err = nil

	return nil
}

func (s *mqttString) Get() (string, error) {
	if err := s.err; err != nil {
		return "", err
	}

	return s.String.Get()
}

func (s *mqttString) Close() error {
	if s.conn == nil {
		return nil
	}

	s.conn.Unsubscribe(s.topic)
	s.conn = nil

	return nil
}
