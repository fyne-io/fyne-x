package main

import (
	"bluetoothFyne/bluetooth"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func runClientReturnError(w fyne.Window) {
	entry, label := widget.NewEntry(), widget.NewLabel("put mac address")
	w.SetContent(container.NewVBox(label, entry,
		widget.NewButton("submit", func() {
			// TODO check mac address
			adapter, err := bluetooth.NewBluetoothDefaultAdapter()
			if err != nil {
				return
			}
			defer fmt.Println(adapter.Close())
			connectTo(adapter, entry.Text)
		})))
}

func connectTo(adapter bluetooth.Adapter, mac string) {
	socket, err := adapter.ConnectAsClientToServer(mac)
	if err != nil {
		return
	}
	defer fmt.Println(socket.Close())
	bytes, er := json.Marshal(msg{}) // defined size of msg
	if er != nil {
		return
	}
	expected := len(bytes)
	for i := uint64(0); i < N; i++ {
		m := msg{}
		n, er := socket.Read(bytes)
		if er != nil || n != expected {
			fmt.Println("1. ", er, ", ", n)
			continue
		}
		er = json.Unmarshal(bytes, &m)
		if er != nil {
			fmt.Println("2. ", er)
			continue
		}
		fmt.Println("3. ", m)
		if m.end {
			return
		}
	}
	// TODO meaningfull logic
}
