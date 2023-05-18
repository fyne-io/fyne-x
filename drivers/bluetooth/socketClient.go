//go:build android

package bluetooth

/*
#include <jni>
*/
import "C"

type BluetoothSocket struct {
	self          C.jobject
	name, address string
}
