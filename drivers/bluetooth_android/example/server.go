package main

import (
	"bluetoothFyne/bluetooth_android" // testing project
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

func runServerReturnError(window fyne.Window) {
	_, err := getAdapter()
	if err != nil {

	}

}

func getAdapter() (*bluetooth_android.BluetoothAdapter, error) {
	var err1 error = nil
	var bluetoothAdapter *bluetooth_android.BluetoothAdapter
	err0 := driver.RunNative(func(context interface{}) error {
		if androidContext, ok := context.(driver.AndroidContext); ok {
			adapter, err := bluetooth_android.NewBluetoothDefaultAdapter(androidContext.Env)
			if err != nil {
				err1 = err
				return nil
			}
			bluetoothAdapter = adapter
		}
		return nil
	})
	return bluetoothAdapter, errors.Join(err0, err1)
}
