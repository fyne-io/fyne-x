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

jobject createBluetoothServer(uintptr_t env, jobject bluetoothAdapter, char** errorMsg) {
    JNIEnv *envPtr = (JNIEnv*)env;

    // Find class BluetoothServerSocket
    jclass serverSocketClass = (*envPtr)->FindClass(envPtr, "android/bluetooth/BluetoothServerSocket");
    if (serverSocketClass == NULL) {
        copyToError( "Failed to find class BluetoothServerSocket", errorMsg);
        return NULL;
    }

    // Get method listenUsingRfcommWithServiceRecord
    jmethodID listenMethod = (*envPtr)->GetMethodID(envPtr, serverSocketClass, "listenUsingRfcommWithServiceRecord", "(Ljava/lang/String;Ljava/util/UUID;)Landroid/bluetooth/BluetoothServerSocket;");
    if (listenMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError("Failed to get method listenUsingRfcommWithServiceRecord", errorMsg);
        return NULL;
    }

    // Create UUID object
    jclass uuidClass = (*envPtr)->FindClass(envPtr, "java/util/UUID");
    if (uuidClass == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError("Failed to find UUID class", errorMsg);
        return NULL;
    }

	// create UUID object
    jmethodID uuidConstructor = (*envPtr)->GetMethodID(envPtr, uuidClass, "<init>", "(JJ)V");
    if (uuidClass == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, uuidClass);
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError("Failed to find uuid constructor", errorMsg);
        return NULL;
    }
    jobject uuidObj = (*envPtr)->NewObject(envPtr, uuidClass, uuidConstructor, (jlong) 0, (jlong) 0);
    if (uuidObj == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, uuidClass);
        copyToError("Failed to create UUID object", errorMsg);
        return NULL;
    }

    // Call listenUsingRfcommWithServiceRecord method
    jobject serverSocket = (*envPtr)->CallObjectMethod(envPtr, bluetoothAdapter, listenMethod, NULL, uuidObj);
    if (serverSocket == NULL) {
        copyToError( "Failed to call listenUsingRfcommWithServiceRecord method", errorMsg);
    }

	// free memory
    (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
    (*envPtr)->DeleteLocalRef(envPtr, uuidClass);
    (*envPtr)->DeleteLocalRef(envPtr, uuidObj);
    return serverSocket;
}

char* getAddress(uintptr_t env, jobject serverSocket, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get BluetoothServerSocket class
    jclass serverSocketClass = (*envPtr)->GetObjectClass(envPtr, serverSocket);
    if (serverSocketClass == NULL) {
        copyToError( "Failed get serverSocket class", errorMsg);
        return NULL;
    }

    // Get the getRemoteDevice method
    jmethodID getRemoteDeviceMethod = (*envPtr)->GetMethodID(envPtr, serverSocketClass, "getRemoteDevice", "()Landroid/bluetooth/BluetoothDevice;");
    if (getRemoteDeviceMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError("Failed to get getRemoteDevice method", errorMsg);
        return NULL;
    }

    // Call the getRemoteDevice method
    jobject bluetoothDevice = (*envPtr)->CallObjectMethod(envPtr, serverSocket, getRemoteDeviceMethod);
    if (bluetoothDevice == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError( "Failed to call the getRemoteDevice method", errorMsg);
        return NULL;
    }

    // Get the getAddress method
    jmethodID getAddressMethod = (*envPtr)->GetMethodID(envPtr, (*envPtr)->GetObjectClass(envPtr, bluetoothDevice), "getAddress", "()Ljava/lang/String;");
    if (getAddressMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDevice);
        copyToError("Failed to get getAddress method", errorMsg);
        return NULL;
    }

    // Call the getAddress method
    jstring addressString = (jstring)(*envPtr)->CallObjectMethod(envPtr, bluetoothDevice, getAddressMethod);
    if (addressString == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothDevice);
        copyToError( "Address string is NULL", errorMsg);
        return NULL;
    }

    // Convert jstring to char*
    const char* addressChars = (*envPtr)->GetStringUTFChars(envPtr, addressString, NULL);
    char* address = (char*)malloc((strlen(addressChars) + 1) * sizeof(char));
    strcpy(address, addressChars);

    // Release memory
    (*envPtr)->ReleaseStringUTFChars(envPtr, addressString, addressChars);
    (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothDevice);
    (*envPtr)->DeleteLocalRef(envPtr, addressString);

    return address;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type ServerSocket struct {
	self    C.jobject
	address string
}

// GetBluetoothServerSocket returns Socket which is listening on bluetooth
func (b *Adapter) GetBluetoothServerSocket() (bs *ServerSocket, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		socket := C.createBluetoothServer(C.uintptr_t(env), b.self, &errMsgC)
		if errMsgC != nil {
			err := errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			bs, e = nil, err
			return nil
		}
		bs = &ServerSocket{self: socket}
		return nil
	})
	e = errors.Join(e, err, bs.FetchAddress())
	return
}

// FetchAddress it is usefully if GetAddress return empty string, it try to set internal address
func (b *ServerSocket) FetchAddress() (e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		nameC := C.getAddress(C.uintptr_t(env), b.self, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		b.address = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return errors.Join(e, err)
}

// GetAddress returns address of server
func (b *ServerSocket) GetAddress() string {
	return b.address
}
