//go:build unix && !android

package bluetooth

import "C"
import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// adapterUnix
type adapterUnix struct{}

// socketUnix
type socketUnix struct {
	fd int
}

// serverSocketUnix
type serverSocketUnix struct {
	fd int
}

// NewBluetoothDefaultAdapter get Bluetooth adapter
func NewBluetoothDefaultAdapter() (b Adapter, e error) {
	return &adapterUnix{}, nil
}

func (a *adapterUnix) FetchAddress() (string, error) {
	return os.Hostname()
}

func (a *adapterUnix) Close() error {
	return nil
}

func (a *adapterUnix) GetBluetoothServerSocket() (ServerSocket, error) {
	fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	if err != nil {
		return nil, err
	}
	addr := &unix.SockaddrRFCOMM{Channel: 1}
	err = unix.Bind(fd, addr)
	if err != nil {
		return nil, errors.Join(err, unix.Close(fd))
	}
	err = unix.Listen(fd, 10)
	if err != nil {
		return nil, errors.Join(err, unix.Close(fd))
	}
	return &serverSocketUnix{fd: fd}, nil
}

func (a *adapterUnix) ConnectAsClientToServer(address string) (Socket, error) {
	mac := a.str2ba(address)
	fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	if err != nil {
		return nil, err
	}
	addr := &unix.SockaddrRFCOMM{Addr: mac, Channel: 1}
	err = unix.Connect(fd, addr)
	if err != nil {
		return nil, errors.Join(err, unix.Close(fd))
	}
	return &socketUnix{fd: fd}, nil
}

func (a *adapterUnix) str2ba(s string) [6]byte {
	p := strings.Split(s, ":")
	var b [6]byte
	for i, tmp := range p {
		u, _ := strconv.ParseUint(tmp, 16, 8)
		b[len(b)-1-i] = byte(u)
	}
	return b
}

func (s *serverSocketUnix) FetchStringData() (string, error) {
	getpeername, err := unix.Getpeername(s.fd)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(getpeername), err
}

func (s *serverSocketUnix) Close() error {
	return unix.Close(s.fd)
}

func (s *serverSocketUnix) Accept() (Socket, error) {
	connFd, _, err := unix.Accept(s.fd)
	if err != nil {
		return nil, err
	}
	return &socketUnix{fd: connFd}, nil
}

func (s *socketUnix) FetchStringData() (string, error) {
	getpeername, err := unix.Getpeername(s.fd)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(getpeername), err
}

func (s *socketUnix) Close() error {
	return unix.Close(s.fd)
}

func (s *socketUnix) Read(bytes []byte) (int, error) {
	return unix.Read(s.fd, bytes)
}

func (s *socketUnix) Write(bytes []byte) (int, error) {
	return unix.Write(s.fd, bytes)
}
