//go:build !unix && !android && windows

package bluetooth

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"strconv"
	"strings"
	"syscall"
)

// adapter
type adapterWin struct{}

// socket
type socketWin struct {
	address windows.Sockaddr
	fd      windows.Handle
}

// serverSocket
type serverSocketWin struct {
	fd windows.Handle
}

// NewBluetoothDefaultAdapter get Bluetooth adapter
func NewBluetoothDefaultAdapter() (b Adapter, e error) {
	return &adapterWin{}, nil
}

func (a *adapterWin) GetAddress() (string, error) {
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

func (a *adapterWin) Close() error {
	return nil
}

func (a *adapterWin) GetBluetoothServerSocket() (ServerSocket, error) {
	var d syscall.WSAData
	err := syscall.WSAStartup(uint32(0x202), &d)
	if err != nil {
		return nil, err
	}

	fd, err := windows.Socket(windows.AF_BTH, windows.SOCK_STREAM, windows.BTHPROTO_RFCOMM)
	if err != nil {
		return nil, err
	}
	s := &windows.SockaddrBth{
		BtAddr: 0,
		Port:   0,
	}
	err = windows.Bind(fd, s)
	if err != nil {
		_ = windows.Close(fd)
		return nil, err
	}
	return &serverSocketWin{fd: fd}, nil
}

func (a *adapterWin) ConnectAsClientToServer(address string) (Socket, error) {
	var d syscall.WSAData
	err := syscall.WSAStartup(uint32(0x202), &d)
	if err != nil {
		return nil, err
	}

	fd, err := windows.Socket(windows.AF_BTH, windows.SOCK_STREAM, windows.BTHPROTO_RFCOMM)
	if err != nil {
		return nil, err
	}

	addressUint64, err := addressToUint64(address)
	if err != nil {
		return nil, err
	}
	s := &windows.SockaddrBth{
		BtAddr: addressUint64,
		Port:   6,
	}
	if err = windows.Connect(fd, s); err != nil {
		return nil, err
	}
	return &socketWin{fd: fd, address: s}, nil
}

// addressToUint64 converts MAC address string to uint64
func addressToUint64(address string) (uint64, error) {
	addressParts := strings.Split(address, ":")
	addressPartsLength := len(addressParts)
	var result uint64
	for i, tmp := range addressParts {
		u, err := strconv.ParseUint(tmp, 16, 8)
		if err != nil {
			return 0, err
		}
		push := 8 * (addressPartsLength - 1 - i)
		result += u << push
	}
	return result, nil
}

func (s *socketWin) Close() error {
	return windows.Close(s.fd)
}

func (s *socketWin) Read(bytes []byte) (int, error) {
	flags := uint32(0)
	buf := windows.WSABuf{Len: uint32(len(bytes)), Buf: &bytes[0]}
	receiver := uint32(0)
	err := windows.WSARecv(s.fd, &buf, 1, &receiver, &flags, nil, nil)
	if err != nil {
		return 0, err
	}
	return int(receiver), nil
}

func (s *socketWin) Write(bytes []byte) (int, error) {
	buf := &windows.WSABuf{
		Len: uint32(len(bytes)),
	}
	if len(bytes) > 0 {
		buf.Buf = &bytes[0]
	}
	var numOfBytes uint32
	err := windows.WSASend(s.fd, buf, 1, &numOfBytes, 0, nil, nil)
	return int(numOfBytes), err
}

func (s *serverSocketWin) Close() error {
	return windows.Close(s.fd)
}

func (s *serverSocketWin) Accept() (Socket, error) {
	accept, addr, err := windows.Accept(s.fd)
	if err != nil {
		return nil, err
	}
	return &socketWin{
		address: addr,
		fd:      accept,
	}, nil
}
func (s *socketWin) StringData() (string, error) {
	adr, ok := s.address.(*windows.SockaddrBth)
	if !ok {
		return "", errors.New("unexpected error casting *windows.SockaddrBth")
	}
	return fmt.Sprint("mac: ", uint64AddressToStr(adr.BtAddr)), nil
}

func uint64AddressToStr(macAddress uint64) string {
	macString := fmt.Sprintf("%012X", macAddress)
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		macString[0:2], macString[2:4], macString[4:6],
		macString[6:8], macString[8:10], macString[10:12])

}
