//go:build android

package bluetooth

var runOnJVM func(func(vm, env, ctx uintptr) error) error = nil

func SetVMFunc(fn func(func(vm, env, ctx uintptr) error) error) {
	if fn == nil {
		panic("you must set to concrete func provide by frame work, for example ")
	}
	runOnJVM = fn
}
