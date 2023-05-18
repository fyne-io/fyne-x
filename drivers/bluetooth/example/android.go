//go:build android

package main

import (
	"bluetoothFyne/bluetooth"
	"errors"
	"fyne.io/fyne/v2/driver"
)

func compatibleWithGOMobile(fn func(vm, env, ctx uintptr) error) error {
	return driver.RunNative(func(context interface{}) error {
		if androidContext, ok := context.(driver.AndroidContext); ok {
			return fn(androidContext.VM, androidContext.Env, androidContext.Ctx)
		}
		return errors.New("no context android")
	})
}

func init() {
	bluetooth.SetVMFunc(compatibleWithGOMobile)
}
