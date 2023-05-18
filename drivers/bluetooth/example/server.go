package main

import (
	"bluetoothFyne/bluetooth" // testing project
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
)

type msg struct {
	x, y, z uint64
}

// example of server side
func runServerReturnError(window fyne.Window) {
	adapter, err := bluetooth.NewBluetoothDefaultAdapter()
	defer adapter.Close() // error can means only path of fail
	if err != nil {
		return
	}
	serverAutomatic(adapter) // autonatic many
	serverManual(adapter)    // manual one
	buildMainMenu(window)
}

func comunicate(readWriter *bluetooth.ReadWriterBluetooth, socketInfo *bluetooth.BluetoothSocket) {
	m := msg{
		x: 10,
		y: 10000,
		z: 5,
	}
	for i := uint64(0); i < 100000000000000; i++ {
		m.x = (m.x + i + m.y + m.z) % 300
		m.y = (m.y + i + m.x + m.z) % 500
		m.z = (m.z + i + m.x + m.y) % 700
		bytes, er := json.Marshal(m) // defined size of msg
		if er == nil {
			_, err := readWriter.Write(bytes)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		fmt.Println(er)
	}

}

func serverAutomatic(adapter *bluetooth.BluetoothAdapter) {
	err := adapter.Server(comunicate)
	if err != nil {
		fmt.Println(err)
	}
}

func serverManual(adapter *bluetooth.BluetoothAdapter) {
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
	comunicate(readWriter, con)
}
