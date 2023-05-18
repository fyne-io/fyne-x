//go:build android

package bluetooth

/*
#include <jni.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static inline void copyToError(const char* errorMsg, char** error) {
    *error = (char*) malloc((strlen(errorMsg) + 1) * sizeof(char));
    strcpy(*error, errorMsg);
}

jobject connectToBluetoothServer(uintptr_t env, jobject bluetoothAdapter, const char* serverAddress, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*) env;

    // Get the BluetoothDevice class
    jclass bluetoothDeviceClass = (*envPtr)->FindClass(envPtr, "android/bluetooth/BluetoothDevice");
    if (bluetoothDeviceClass == NULL) {
        copyToError("Failed to find BluetoothDevice class", errorMsg);
        return NULL;
    }

    // Get the getRemoteDevice method from BluetoothAdapter
    jclass bluetoothAdapterClass = (*envPtr)->GetObjectClass(envPtr, bluetoothAdapter);
    if (bluetoothAdapterClass == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDeviceClass);
        copyToError("Failed to get BluetoothAdapter class", errorMsg);
        return NULL;
    }

    jmethodID getRemoteDeviceMethod = (*envPtr)->GetMethodID(envPtr, bluetoothAdapterClass, "getRemoteDevice", "(Ljava/lang/String;)Landroid/bluetooth/BluetoothDevice;");
    if (getRemoteDeviceMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDeviceClass);
        copyToError("Failed to get getRemoteDevice method", errorMsg);
        return NULL;
    }

    // Get the BluetoothDevice object
    jstring serverAddressString = (*envPtr)->NewStringUTF(envPtr, serverAddress);
    jobject bluetoothDevice = (*envPtr)->CallObjectMethod(envPtr, bluetoothAdapter, getRemoteDeviceMethod, serverAddressString);
    if (bluetoothDevice == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDeviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, serverAddressString);
        copyToError("Failed to get the BluetoothDevice object", errorMsg);
        return NULL;
    }

    // Get the createRfcommSocket method from BluetoothDevice
    jmethodID createRfcommSocketMethod = (*envPtr)->GetMethodID(envPtr, bluetoothDeviceClass, "createRfcommSocket", "(I)Landroid/bluetooth/BluetoothSocket;");
    if (createRfcommSocketMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDeviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, serverAddressString);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDevice);
        copyToError("Failed to get createRfcommSocket method", errorMsg);
        return NULL;
    }

    // Create the BluetoothSocket object
    jobject bluetoothSocket = (*envPtr)->CallObjectMethod(envPtr, bluetoothDevice, createRfcommSocketMethod, 1);
    if (bluetoothSocket == NULL) {
        copyToError("Failed: bluetoothSocket == NULL", errorMsg);
    }

    // Clean up local references
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothDeviceClass);
    (*envPtr)->DeleteLocalRef(envPtr, serverAddressString);
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothDevice);

    return bluetoothSocket;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// ConnectAsClientToServer take address of server and return conection in BluetoothSocket, WARNING: error can means only path of fail put defer before error handling
func (b *BluetoothServerSocket) ConnectAsClientToServer(address string) (bs *BluetoothSocket, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		caddress := C.CString(address)
		defer C.free(unsafe.Pointer(caddress))
		sock := C.connectToBluetoothServer(C.uintptr_t(env), b.self, caddress, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		bs = &BluetoothSocket{self: sock}
		return nil
	})
	var e1, e2 error
	if bs != nil {
		e1 = bs.FetchAddress()
		e2 = bs.FetchName()
	}
	return bs, errors.Join(err, e, e1, e2)
}
