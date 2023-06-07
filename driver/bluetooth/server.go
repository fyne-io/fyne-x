package bluetooth

import (
	"errors"
)

type Handle func(ReadWriter, Socket)

// Server creates server and listen to clients
func Server(a Adapter, stop chan struct{}, fn Handle) (resErr error) {
	type errMsg struct {
		err error
		con Socket
	}
	socket, err := a.GetBluetoothServerSocket()
	if err != nil {
		return err
	}
	defer func() { resErr = errors.Join(resErr, socket.Close()) }()
	chMsg := make(chan errMsg)
	for {
		select {
		case <-stop:
			return
		case msg := <-chMsg:
			resErr = msg.err
			if msg.con != nil {
				resErr = errors.Join(resErr, msg.con.Close())
			}
			return
		default:
		}
		con, err0 := socket.Accept()
		if err0 != nil {
			return err0
		}
		go func() {
			readWriter, e := con.GetReadWriter()
			if e != nil {
				chMsg <- errMsg{err: e, con: con}
				return
			}
			fn(readWriter, con) // handle
			e = readWriter.Close()
			if e != nil {
				chMsg <- errMsg{err: e, con: con}
				return
			}
			e = con.Close()
			if e != nil {
				chMsg <- errMsg{err: e, con: nil}
			}
		}()
	}
}
