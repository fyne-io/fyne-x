//go:build unix && !android

package bluetooth

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"strconv"
	"strings"
	"syscall"
)

// adapter
type adapterUnix struct{}

// socket
type socketUnix struct {
	fd      int
	address unix.Sockaddr
}

// serverSocket
type serverSocketUnix struct {
	fd int
}

// NewBluetoothDefaultAdapter get Bluetooth adapter
func NewBluetoothDefaultAdapter() (b Adapter, e error) {
	return &adapterUnix{}, nil
}

func (a *adapterUnix) GetAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				return i.HardwareAddr.String(), nil
			}
		}
	}
	return "", errors.New("not found")
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
	return &socketUnix{fd: fd, address: addr}, nil
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

func (s *serverSocketUnix) Close() error {
	return unix.Close(s.fd)
}

func (s *serverSocketUnix) Accept() (Socket, error) {
	connFd, address, err := unix.Accept(s.fd)
	if err != nil {
		return nil, err
	}
	return &socketUnix{fd: connFd, address: address}, nil
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

func (s *socketUnix) StringData() (string, error) {
	adr, ok := s.address.(*unix.SockaddrRFCOMM)
	if !ok {
		return "", errors.New("unexpected error casting *windows.SockaddrRFCOMM")
	}
	return fmt.Sprint("mac: ", uint8AddressToStr(adr.Addr)), nil
}

func uint8AddressToStr(s [6]uint8) string {
	return fmt.Sprint(s[0], ":", s[1], ":", s[2], ":", s[3], ":", s[4], ":", s[5], ":")
}
