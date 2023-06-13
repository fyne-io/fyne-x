//go:build !unix && !android && windows

package bluetooth

import "errors"

// adapter
type adapterWin struct{}

// socket
type socketWin struct{}

// serverSocket
type serverSocketWin struct{}

// NewBluetoothDefaultAdapter get Bluetooth adapter
func NewBluetoothDefaultAdapter() (b Adapter, e error) {
	return nil, errors.New("not implemented")
}

func (a *adapterWin) FetchAddress() (string, error) {
	return "", errors.New("not implemented")
}

func (a *adapterWin) Close() error {
	return errors.New("not implemented")
}

func (a *adapterWin) GetBluetoothServerSocket() (ServerSocket, error) {
	return nil, errors.New("not implemented")
}

func (a *adapterWin) ConnectAsClientToServer(s string) (Socket, error) {
	return nil, errors.New("not implemented")
}

func (s *socketWin) FetchStringData() (string, error) {
	return "", errors.New("not implemented")
}

func (s *socketWin) Close() error {
	return errors.New("not implemented")
}

func (s *socketWin) Read(bytes []byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (s *socketWin) Write(bytes []byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (s *serverSocketWin) FetchStringData() (string, error) {
	return "", errors.New("not implemented")
}

func (s *serverSocketWin) Close() error {
	return errors.New("not implemented")
}

func (s *serverSocketWin) Accept() (Socket, error) {
	return nil, errors.New("not implemented")
}
