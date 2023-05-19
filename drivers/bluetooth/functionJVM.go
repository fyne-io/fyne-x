//go:build android

package bluetooth

import (
	"errors"
	"fyne.io/fyne/v2/driver"
)

var runOnJVM = originForm

func originForm(fn func(vm, env, ctx uintptr) error) error {
	return driver.RunNative(func(context interface{}) error {
		if androidContext, ok := context.(driver.AndroidContext); ok {
			return fn(androidContext.VM, androidContext.Env, androidContext.Ctx)
		}
		return errors.New("no context android")
	})
}
