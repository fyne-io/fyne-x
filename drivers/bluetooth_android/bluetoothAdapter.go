//go:build android

package bluetooth_android

/*
#include <jni.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static inline void copyToError(char* errorMsgv, char** error){
    *error = (char*)malloc((strlen(errorMsg) + 1) * sizeof(char));
    strcpy(*error, errorMsg);
}

jobject getBluetoothAdapter(uintptr_t env, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get the BluetoothAdapter class
    jclass bluetoothAdapterClass = (*envPtr)->FindClass(envPtr, "android/bluetooth/BluetoothAdapter");
    if (bluetoothAdapterClass == NULL) {
        copyToError("Failed to find BluetoothAdapter class", errorMsg);
        return NULL;
    }

    // Get the getDefaultAdapter static method
    jmethodID getDefaultAdapterMethod = (*envPtr)->GetStaticMethodID(envPtr, bluetoothAdapterClass, "getDefaultAdapter", "()Landroid/bluetooth/BluetoothAdapter;");
    if (getDefaultAdapterMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        copyToError("Failed to find getDefaultAdapter method", errorMsg);
        return NULL;
    }

    // Get the BluetoothAdapter object
    jobject bluetoothAdapter = (*envPtr)->CallStaticObjectMethod(envPtr, bluetoothAdapterClass, getDefaultAdapterMethod);
    if (bluetoothAdapter == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        copyToError("Failed to invoke getDefaultAdapter method", errorMsg);
        return NULL;
    }

    // Release memory
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
    return bluetoothAdapter;
}

char* getBluetoothName(uintptr_t env, jobject bluetoothAdapter, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get BluetoothAdapter class
    jclass bluetoothAdapterClass = (*envPtr)->GetObjectClass(envPtr, bluetoothAdapter);
    if (bluetoothAdapterClass == NULL) {
        copyToError("Failed to find BluetoothAdapter class", errorMsg);
        return NULL;
    }

    // Get the getName method
    jmethodID getNameMethod = (*envPtr)->GetMethodID(envPtr, bluetoothAdapterClass, "getName", "()Ljava/lang/String;");
    // If the method is not available, return NULL
    if (getNameMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        copyToError("Failed to get getName method", errorMsg);
        return NULL;
    }

    // Call the getName method
    jstring bluetoothName = (jstring)(*envPtr)->CallObjectMethod(envPtr, bluetoothAdapter, getNameMethod);
    if (bluetoothName == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        copyToError("Bluetooth name is NULL", errorMsg);
        return NULL;
    }

    // Convert jstring to char*
    const char* nameChars = (*envPtr)->GetStringUTFChars(envPtr, bluetoothName, NULL);
    if (bluetoothName == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        copyToError("Failed convert jstring to char*", errorMsg);
        return NULL;
    }
    size_t length = strlen(nameChars);
    char* heapString = (char*)malloc((length + 1) * sizeof(char));
    strcpy(heapString, nameChars);

    // Release memory
    (*envPtr)->ReleaseStringUTFChars(envPtr, bluetoothName, nameChars);
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothName);

    return heapString;
}

void freeBluetoothAdapter(uintptr_t env, jobject bluetoothAdapter, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    if (bluetoothAdapter != NULL) {
        // Get BluetoothAdapter class
        jclass bluetoothAdapterClass = (*envPtr)->GetObjectClass(envPtr, bluetoothAdapter);

        // Get the close method
        jmethodID closeMethod = (*envPtr)->GetMethodID(envPtr, bluetoothAdapterClass, "close", "()V");

        // If the method is not available, set error and return
        if (closeMethod == NULL) {
            copyToError( "Failed to get close method", errorMsg);
            return;
        }

        // Call the close method
        (*envPtr)->CallVoidMethod(envPtr, bluetoothAdapter, closeMethod);

        // Delete local references
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapter);
    }
}
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type BluetoothAdapter struct {
	self C.jobject
	name string
	stop chan struct{}
}

func NewBluetoothDefaultAdapter(env uintptr) (*BluetoothAdapter, error) {
	var errMsgC *C.char

	adapter := C.getBluetoothAdapter(C.uintptr_t(env), &errMsgC)
	if errMsgC != nil {
		err := errors.New(C.GoString(errMsgC))
		C.free(unsafe.Pointer(errMsgC))
		return nil, err
	}
	result := &BluetoothAdapter{self: adapter}
	nameC := C.getBluetoothName(C.uintptr_t(env), adapter, &errMsgC)
	if errMsgC != nil {
		err := errors.New(C.GoString(errMsgC))
		C.free(unsafe.Pointer(errMsgC))
		fmt.Println(err)
		return result, nil
	}
	result.name = C.GoString(nameC)
	C.free(unsafe.Pointer(nameC))
	return result, nil
}

// GetName returns name of device
func (b *BluetoothAdapter) GetName() string {
	return b.name
}

// Close clean bluetooth adapter from memory
func (b *BluetoothAdapter) Close(env uintptr) {
	var errMsgC *C.char
	C.freeBluetoothAdapter(C.uintptr_t(env), b.self, &errMsgC)
	if errMsgC != nil {
		err := errors.New(C.GoString(errMsgC))
		C.free(unsafe.Pointer(errMsgC))
		fmt.Println(err)
	}
}

// Server creates server and listen to client lim
func (b *BluetoothAdapter) Server(env uintptr, howMany uint64, fn func(*BluetoothSocket)) error {
	socket, err := b.GetBluetoothServerSocket(env)
	if err != nil {
		return err
	}
	defer func(socket *BluetoothServerSocket, env uintptr) {
		err := socket.Close(env)
		if err != nil {
			fmt.Println(err)
		}
	}(socket, env)
	for i := uint64(0); i < howMany; i++ {
		select {
		case <-b.stop:
			return nil
		default:
		}
		con, err0 := socket.Accept(env)
		if err0 != nil {
			return err0
		}
		go fn(con)
	}
	return nil
}
