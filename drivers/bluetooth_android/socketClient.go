//go:build android

package bluetooth_android

/*
#include <jni>
*/
import "C"

type BluetoothSocket struct {
	self          C.jobject
	name, address string
}
