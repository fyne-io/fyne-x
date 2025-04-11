package bluetooth

import (
	"errors"
)

type Handle func(Socket)

// Server creates server and listen to clients
func Server(a Adapter, stop chan struct{}, fn Handle) (resErr error) {
	socket, err := a.GetBluetoothServerSocket()
	if err != nil {
		return err
	}
	defer func() { resErr = errors.Join(resErr, socket.Close()) }()
	chErr := make(chan error)
	for {
		select {
		case <-stop:
			return
		case e := <-chErr:
			resErr = e
			return
		default:
		}
		con, err0 := socket.Accept()
		if err0 != nil {
			resErr = err0
			return
		}
		go func() {
			fn(con)
			e := con.Close()
			if e != nil {
				chErr <- e
			}
		}()
	}
}
