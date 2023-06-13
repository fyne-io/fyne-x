//go:build !unix && !android && !windows

package bluetooth

import "errors"

// adapter
type adapterOther struct{}

// socket
type socketOther struct{}

// serverSocket
type serverSocketOther struct{}

// NewBluetoothDefaultAdapter get Bluetooth adapter
func NewBluetoothDefaultAdapter() (b Adapter, e error) {
	return nil, errors.New("not implemented")
}

func (a *adapterOther) FetchAddress() (string, error) {
	return "", errors.New("not implemented")
}

func (a *adapterOther) Close() error {
	return errors.New("not implemented")
}

func (a *adapterOther) GetBluetoothServerSocket() (ServerSocket, error) {
	return nil, errors.New("not implemented")
}

func (a *adapterOther) ConnectAsClientToServer(s string) (Socket, error) {
	return nil, errors.New("not implemented")
}

func (s *socketOther) FetchStringData() (string, error) {
	return "", errors.New("not implemented")
}

func (s *socketOther) Close() error {
	return errors.New("not implemented")
}

func (s *socketOther) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (s *socketOther) Write([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (s *serverSocketOther) FetchStringData() (string, error) {
	return "", errors.New("not implemented")
}

func (s *serverSocketOther) Close() error {
	return errors.New("not implemented")
}

func (s *serverSocketOther) Accept() (Socket, error) {
	return nil, errors.New("not implemented")
}
