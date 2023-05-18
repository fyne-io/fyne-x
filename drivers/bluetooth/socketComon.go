//go:build android

package bluetooth

/*
#include <jni>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static inline void copyToError(char* errorMsgv, char** error){
    *error = (char*)malloc((strlen(errorMsg) + 1) * sizeof(char));
    strcpy(*error, errorMsg);
}

void closeBluetoothServerSocket(uintptr_t env, jobject serverSocket, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get the BluetoothServerSocket class
    jclass serverSocketClass = (*envPtr)->GetObjectClass(envPtr, serverSocket);

    // Get the close method
    jmethodID closeMethod = (*envPtr)->GetMethodID(envPtr, serverSocketClass, "close", "()V");

    // If the close method is available, invoke it
    if (closeMethod != NULL) {
        (*envPtr)->CallVoidMethod(envPtr, serverSocket, closeMethod);
    } else {
        copyToError("Failed to find close method", errorMsg);
    }

    // Release memory
    (*envPtr)->DeleteLocalRef(envPtr, serverSocket);
    (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

func (b *BluetoothServerSocket) Close() error {
	var errMsgC *C.char
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		C.closeBluetoothServerSocket(C.uintptr_t(env), b.self, &errMsgC)
		return nil
	})

	if errMsgC != nil {
		err = errors.Join(errors.New(C.GoString(errMsgC)), err)
		C.free(unsafe.Pointer(errMsgC))
		return err
	}
	return nil
}

func (b *BluetoothSocket) Close() error {
	var errMsgC *C.char
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		C.closeBluetoothServerSocket(C.uintptr_t(env), b.self, &errMsgC)
		return nil
	})

	if errMsgC != nil {
		err = errors.Join(errors.New(C.GoString(errMsgC)), err)
		C.free(unsafe.Pointer(errMsgC))
		return err
	}
	return nil
}
