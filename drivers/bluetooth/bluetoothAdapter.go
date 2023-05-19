//go:build android

package bluetooth

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
    if (bluetoothAdapter != NULL) {
        JNIEnv* envPtr = (JNIEnv*)env;
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

// Adapter has references android.bluetooth.BluetoothAdapter in java
type Adapter struct {
	self C.jobject
	name string
	stop chan struct{}
}

// StopServer send signal to server to end
func (b *Adapter) StopServer() {
	go func() {
		b.stop <- struct{}{}
	}()
}

// NewBluetoothDefaultAdapter get Bluetooth adapter, WARNING: error can means only path of fail put defer before error handling
func NewBluetoothDefaultAdapter() (b *Adapter, e error) {
	if runOnJVM == nil {
		return nil, errors.New("you must on android call SetVMFunc before")
	}
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		adapter := C.getBluetoothAdapter(C.uintptr_t(env), &errMsgC)
		if errMsgC != nil {
			err := errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			e = err
			return nil
		}
		b = &Adapter{self: adapter}
		nameC := C.getBluetoothName(C.uintptr_t(env), adapter, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		b.name = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	e = errors.Join(e, err)
	return
}

// GetName returns name of device
func (b *Adapter) GetName() string {
	return b.name
}

// Close clean bluetooth adapter from memory
func (b *Adapter) Close() {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		C.freeBluetoothAdapter(C.uintptr_t(env), b.self, &errMsgC)
		if errMsgC != nil {
			err := errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
}

// Server creates server and listen to clients
func (b *Adapter) Server(fn Handle) error {
	socket, err := b.GetBluetoothServerSocket()
	defer fmt.Println(socket.Close())
	if err != nil {
		return err
	}
	for {
		select {
		case <-b.stop:
			return nil
		default:
		}
		con, err0 := socket.Accept()
		if err0 != nil {
			fmt.Println(con.Close())
			return err0
		}
		go func() {
			defer fmt.Println(con.Close())
			readWriter, er := con.GetReadWriter()
			defer fmt.Println(readWriter.Close())
			if er != nil {
				return
			}
			fn(readWriter, con)
		}()
	}
}
