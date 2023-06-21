package main

import (
	"fmt"
	"fyne.io/fyne/v2/driver/bluetooth"
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
	comunicate2(con)
}

func comunicate2(con bluetooth.Socket) {
	// TODO meaningfull logic
}
