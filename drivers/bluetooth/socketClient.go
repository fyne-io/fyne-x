//go:build android

package bluetooth

/*
#include <jni>
*/
import "C"

type Socket struct {
	self          C.jobject
	name, address string
}
