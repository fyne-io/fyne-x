package main

import (
	"bluetoothFyne/bluetooth"
	"fmt"
)

func serverManual(adapter bluetooth.Adapter) {
	socket, er := adapter.GetBluetoothServerSocket()
	defer fmt.Println(socket.Close())
	if er != nil {
		fmt.Println(er)
		return
	}
	con, er := socket.Accept()
	defer fmt.Println(con.Close())
	if er != nil {
		fmt.Println(er)
		return
	}
	readWriter, er := con.GetReadWriter()
	defer fmt.Println(readWriter.Close())
	if er != nil {
		return
	}
	comunicate2(readWriter, con)
}

func comunicate2(readWriter bluetooth.ReadWriter, con bluetooth.Socket) {

	// TODO meaningfull logic
}
