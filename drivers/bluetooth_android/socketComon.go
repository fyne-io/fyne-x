//go:build android

package bluetooth_android

/*
##include <jni>
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

func (b *BluetoothServerSocket) Close(env uintptr) error {
	var errMsgC *C.char
	C.closeBluetoothServerSocket(C.uintptr_t(env), b.self, &errMsgC)
	if errMsgC != nil {
		err := errors.New(C.GoString(errMsgC))
		C.free(unsafe.Pointer(errMsgC))
		return err
	}
	return nil
}

func (b *BluetoothSocket) Close(env uintptr) error {
	var errMsgC *C.char
	C.closeBluetoothServerSocket(C.uintptr_t(env), b.self, &errMsgC)
	if errMsgC != nil {
		err := errors.New(C.GoString(errMsgC))
		C.free(unsafe.Pointer(errMsgC))
		return err
	}
	return nil
}
