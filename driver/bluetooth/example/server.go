package main

import "C"
import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/bluetooth"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

// example of server side
func runServerReturnError(window fyne.Window) {
	adapter, err := bluetooth.NewBluetoothDefaultAdapter()
	if err != nil {
		return
	}
	defer fmt.Println(adapter.Close())
	serverAutomatic(adapter, window)
}

// handle function
func comunicate(socket bluetooth.Socket) {
	m := msg{
		x: 10,
		y: 10000,
		z: 5,
	}
	b, _ := howMany.Get()
	b++
	_ = howMany.Set(b)
	_ = strHowMany.Set(strconv.Itoa(b))
	for i := uint64(0); i < N; i++ {
		m.x = (m.x + i + m.y + m.z) % 300
		m.y = (m.y + i + m.x + m.z) % 500
		m.z = (m.z + i + m.x + m.y) % 700
		bytes, er := json.Marshal(m) // defined size of msg
		if er == nil {
			_, err := socket.Write(bytes)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		fmt.Println(er)
	}
	// TODO meaningfull logic
}

func serverAutomatic(adapter bluetooth.Adapter, window fyne.Window) {
	stop := make(chan struct{})
	go func() {
		err := bluetooth.Server(adapter, stop, comunicate)
		if err != nil {
			fmt.Println(err)
		}
	}()
	i := 0
	strI := strconv.Itoa(i)
	howMany = binding.BindInt(&i)
	strHowMany = binding.BindString(&strI)
	window.SetContent(container.NewVBox(
		widget.NewLabelWithData(strHowMany),
		widget.NewButton("stop server", func() {
			stop <- struct{}{}
			buildMainMenu(window)
		}),
	))
}
