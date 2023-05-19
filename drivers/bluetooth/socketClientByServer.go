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

char* getClientName(uintptr_t env, jobject clientSocket, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;
    // Get the BluetoothSocket class
    jclass socketClass = (*envPtr)->GetObjectClass(envPtr, clientSocket);
    if (socketClass == NULL) {
        copyToError( "Failed to get getRemoteDevice method", errorMsg);
        return NULL;
    }

    // Get the getRemoteDevice method
    jmethodID remoteDeviceMethod = (*envPtr)->GetMethodID(envPtr, socketClass, "getRemoteDevice", "()Landroid/bluetooth/BluetoothDevice;");
    if (remoteDeviceMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        copyToError("Failed to get getRemoteDevice method", errorMsg);
        return NULL;
    }

    jobject remoteDevice = (*envPtr)->CallObjectMethod(envPtr, clientSocket, remoteDeviceMethod);
    if (remoteDevice == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        copyToError( "Failed to get remote device", errorMsg);
        return NULL;
    }

    // Get the BluetoothDevice class
    jclass deviceClass = (*envPtr)->GetObjectClass(envPtr, remoteDevice);
    if (remoteDevice == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get the BluetoothDevice class", errorMsg);
        return NULL;
    }

    // Get the getName method
    jmethodID getNameMethod = (*envPtr)->GetMethodID(envPtr, deviceClass, "getName", "()Ljava/lang/String;");
    if (getNameMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get getName method", errorMsg);
        return NULL;
    }

    jstring nameString = (*envPtr)->CallObjectMethod(envPtr, remoteDevice, getNameMethod);
    if (nameString == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get name string", errorMsg);
        return NULL;
    }

    // Convert the Java string to C string
    const char* name = (*envPtr)->GetStringUTFChars(envPtr, nameString, NULL);
    if (name == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(env, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        (*envPtr)->ReleaseStringUTFChars(envPtr, nameString, name);
        copyToError( "Failed to convert name string", errorMsg);
        return NULL;
    }

    // Allocate memory for the result
    char* result = (char*)malloc((strlen(name) + 1) * sizeof(char));
    // Copy the name to the result
    strcpy(result, name);

    // Clean up references
    (*envPtr)->DeleteLocalRef(envPtr, socketClass);
    (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
    (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
    (*envPtr)->ReleaseStringUTFChars(envPtr, nameString, name);

	return result;
}

jobject acceptBluetoothClient(uintptr_t env, jobject serverSocket, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get the accept method from BluetoothServerSocket
    jclass serverSocketClass = (*envPtr)->GetObjectClass(envPtr, serverSocket);
    if (serverSocketClass == NULL) {
        copyToError( "Failed to get socket class", errorMsg);
        return NULL;
    }
    jmethodID acceptMethod = (*envPtr)->GetMethodID(envPtr, serverSocketClass, "accept", "()Landroid/bluetooth/BluetoothSocket;");
    if (acceptMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError( "Failed to get accept method", errorMsg);
        return NULL;
    }

    // Call the accept method
    jobject clientSocket = (*envPtr)->CallObjectMethod(envPtr, serverSocket, acceptMethod);
    if (clientSocket == NULL) {
        copyToError("Failed to accept client", errorMsg);
    }
    (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);

    return clientSocket;
}

char* getClientAddress(JNIEnv* env, jobject clientSocket, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;
    // Get the BluetoothSocket class
    jclass socketClass = (*envPtr)->GetObjectClass(envPtr, clientSocket);
    if (serverSocketClass == NULL) {
        copyToError( "Failed to get socket class", errorMsg);
        return NULL;
    }

    // Get the getRemoteDevice method
    jmethodID remoteDeviceMethod = (*envPtr)->GetMethodID(envPtr, socketClass, "getRemoteDevice", "()Landroid/bluetooth/BluetoothDevice;");
    if (remoteDeviceMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        copyToError( "Failed to get getRemoteDevice method", errorMsg);
        return NULL;
    }

    jobject remoteDevice = (*envPtr)->CallObjectMethod(envPtr, clientSocket, remoteDeviceMethod);
    if (remoteDevice == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError("Failed to get remote device", errorMsg);
        return NULL;
    }

    // Get the BluetoothDevice class
    jclass deviceClass = (*envPtr)->GetObjectClass(envPtr, remoteDevice);
    if (getAddressMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get the BluetoothDevice class", errorMsg);
        return NULL;
    }

    // Get the getAddress method
    jmethodID getAddressMethod = (*envPtr)->GetMethodID(envPtr, deviceClass, "getAddress", "()Ljava/lang/String;");
    if (getAddressMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get getAddress method", errorMsg);
        return NULL;
    }

    jstring addressString = (*envPtr)->CallObjectMethod(envPtr, remoteDevice, getAddressMethod);
    if (addressString == NULL) {
        (*envPtr)->DeleteLocalRef(envPtrv, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, serverSocketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError("Failed to get address string", errorMsg);
        return NULL;
    }

    // Convert the Java string to C string
    const char* address = (*envPtr)->GetStringUTFChars(envPtr, addressString, NULL);
    if (address == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
        (*envPtr)->ReleaseStringUTFChars(envPtr, addressString, address);
        copyToError("Failed to convert address string", errorMsg);
        return NULL;
    }

    // Allocate memory for the result
    char* result = (char*)malloc((strlen(address) + 1) * sizeof(char));

    // Copy the address

    strcpy(result, address);

    // Release the Java string and local references
    (*env)->ReleaseStringUTFChars(env, addressString, address);
    (*env)->DeleteLocalRef(env, socketClass);
    (*env)->DeleteLocalRef(env, deviceClass);

    return result;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Accept accepting client and return conection in BluetoothSocket, WARNING: error can means only path of fail put defer before error handling
func (b *ServerSocket) Accept() (bs *Socket, e error) {
	var errMsgC *C.char
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		sock := C.acceptBluetoothClient(C.uintptr_t(env), b.self, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		bs = &Socket{self: sock}
		return nil
	})
	var e1, e2 error

	if bs != nil {
		e1 = bs.FetchName()
		e2 = bs.FetchAddress()
	}
	return bs, errors.Join(err, e, e1, e2)
}

// FetchName it is usefully if GetName return empty string, it try to set internal address
func (b *Socket) FetchName() (e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {

		var errMsgC *C.char
		nameC := C.getClientName(C.uintptr_t(env), b.self, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		b.name = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return errors.Join(err, e)
}

// FetchAddress it is usefully if GetAddress return empty string, it try to set internal address
func (b *Socket) FetchAddress() (e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		nameC := C.getClientAddress(C.uintptr_t(env), b.self, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		b.address = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return errors.Join(err, e)
}

// GetName returns address of client
func (b *Socket) GetName() string {
	return b.name
}

// GetAddress returns address of client
func (b *Socket) GetAddress() string {
	return b.address
}
