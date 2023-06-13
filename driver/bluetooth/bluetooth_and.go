//go:build android

package bluetooth

/*
#include <jni.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static inline void copyToError(char* errorMsg, char** error){
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

char* getClientAddress(uintptr_t env, jobject clientSocket, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;
    // Get the BluetoothSocket class
    jclass socketClass = (*envPtr)->GetObjectClass(envPtr, clientSocket);
    if (socketClass == NULL) {
        copyToError( "Failed to get socket class", errorMsg);
        return NULL;
    }

    // Get the getRemoteDevice method
    jmethodID remoteDeviceMethod = (*envPtr)->GetMethodID(envPtr, socketClass, "getRemoteDevice", "()Landroid/bluetooth/BluetoothDevice;");
    if (remoteDeviceMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        copyToError( "Failed to get getRemoteDevice method", errorMsg);
        return NULL;
    }

    jobject remoteDevice = (*envPtr)->CallObjectMethod(envPtr, clientSocket, remoteDeviceMethod);
    if (remoteDevice == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError("Failed to get remote device", errorMsg);
        return NULL;
    }

    // Get the BluetoothDevice class
    jclass deviceClass = (*envPtr)->GetObjectClass(envPtr, remoteDevice);
    if (deviceClass == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get the BluetoothDevice class", errorMsg);
        return NULL;
    }

    // Get the getAddress method
    jmethodID getAddressMethod = (*envPtr)->GetMethodID(envPtr, deviceClass, "getAddress", "()Ljava/lang/String;");
    if (getAddressMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        (*envPtr)->DeleteLocalRef(envPtr, remoteDevice);
        copyToError( "Failed to get getAddress method", errorMsg);
        return NULL;
    }

    jstring addressString = (*envPtr)->CallObjectMethod(envPtr, remoteDevice, getAddressMethod);
    if (addressString == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, deviceClass);
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
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
    (*envPtr)->ReleaseStringUTFChars(envPtr, addressString, address);
    (*envPtr)->DeleteLocalRef(envPtr, socketClass);
    (*envPtr)->DeleteLocalRef(envPtr, deviceClass);

    return result;
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



jobject getBluetoothOutputStream(uintptr_t env, jobject clientSocket, char** errorMsg) {
    JNIEnv *envPtr = (JNIEnv*)env;
    // Get the BluetoothSocket class
    jclass socketClass = (*envPtr)->GetObjectClass(envPtr, clientSocket);
    if (socketClass == NULL) {
        copyToError( "Failed to get the BluetoothSocket class", errorMsg);
        return NULL;
    }

    // Get the getOutputStream method
    jmethodID getOutputStreamMethod = (*envPtr)->GetMethodID(envPtr, socketClass, "getOutputStream", "()Ljava/io/OutputStream;");
    if (getOutputStreamMethod == NULL) {
        copyToError( "Failed to get getOutputStream method", errorMsg);
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        return NULL;
    }

    // Call the getOutputStream method
    jobject outputStream = (*envPtr)->CallObjectMethod(envPtr, clientSocket, getOutputStreamMethod);
    if (outputStream == NULL) {
        copyToError( "Failed to get OutputStream",  errorMsg);
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
    }

    // Release local references
    (*envPtr)->DeleteLocalRef(envPtr, socketClass);

    return outputStream;
}

jobject getBluetoothInputStream(uintptr_t env, jobject clientSocket, char** errorMsg) {
    JNIEnv *envPtr = (JNIEnv*)env;
    // Get the BluetoothSocket class
    jclass socketClass = (*envPtr)->GetObjectClass(envPtr, clientSocket);
    if (socketClass == NULL) {
        copyToError( "Failed to get the BluetoothSocket class", errorMsg);
        return NULL;
    }

    // Get the getInputStream method
    jmethodID getInputStreamMethod = (*envPtr)->GetMethodID(envPtr, socketClass, "getInputStream", "()Ljava/io/InputStream;");
    if (getInputStreamMethod == NULL) {
        copyToError( "Failed to get getInputStream method", envPtr);
        (*envPtr)->DeleteLocalRef(envPtr, socketClass);
        return NULL;
    }

    // Call the getInputStream method
    jobject inputStream = (*envPtr)->CallObjectMethod(envPtr, clientSocket, getInputStreamMethod);
    if (inputStream == NULL) {
        copyToError( "Failed to get InputStream", errorMsg);
    }

    // Release local references
    (*envPtr)->DeleteLocalRef(envPtr, socketClass);

    return inputStream;
}

void closeStream0(uintptr_t env, jobject stream, char** errorMsg) {
	if (stream == NULL){
	    return;
	    copyToError("empty (null) stream is not able delete", errorMsg);
	}
    JNIEnv *envPtr = (JNIEnv*)env;
	//get stream class
    jclass streamClass = (*envPtr)->GetObjectClass(envPtr, stream);
    if (streamClass == NULL) {
        copyToError( "Failed to get the streamClass class", errorMsg);
        return;
    }

    jmethodID closeMethod = (*envPtr)->GetMethodID(envPtr, streamClass, "close", "()V");
    if (closeMethod != NULL) {
        (*envPtr)->CallVoidMethod(envPtr, stream, closeMethod);
    } else {
        copyToError("Failed to close output stream", errorMsg);
    }
    (*envPtr)->DeleteLocalRef(env, streamClass);
    (*envPtr)->DeleteLocalRef(env, stream);
}

char* readFromInputStream(uintptr_t env, jobject inputStream, int size, int* count, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get the InputStream class
    jclass inputStreamClass = (*envPtr)->GetObjectClass(envPtr, inputStream);
    if (inputStreamClass == NULL) {
        copyToError( "Failed to get the InputStream class", errorMsg);
        return NULL;
    }

    // Get the read method from InputStream
    jmethodID readMethod = (*envPtr)->GetMethodID(envPtr, inputStreamClass, "read", "([B)I");
    if (readMethod == NULL) {
        *count = -1;
        (*envPtr)->DeleteLocalRef(envPtr, inputStreamClass);
        copyToError("Failed to get read method", errorMsg);
        return NULL;
    }

    // Create a byte array for reading data
    jbyteArray byteArray = (*envPtr)->NewByteArray(envPtr, size);
    if (byteArray == NULL) {
        *count = -1;
        (*envPtr)->DeleteLocalRef(envPtr, inputStreamClass);
        copyToError( "Failed to create byte array", errorMsg);
        return NULL;
    }

    // Call the read method
    jint bytesRead = (*envPtr)->CallIntMethod(envPtr, inputStream, readMethod, byteArray);
    if (bytesRead < 0) {
        *count = -1;
        (*envPtr)->DeleteLocalRef(envPtr, inputStreamClass);
        (*envPtr)->DeleteLocalRef(envPtr, byteArray);
        copyToError("Failed to read from input stream", errorMsg);
        return NULL;
    }

    // Copy the bytes from the byte array to a new buffer
    jbyte* byteBuffer = (*envPtr)->GetByteArrayElements(envPtr, byteArray, NULL);
    if (byteArray == NULL) {
        *count = -1;
        (*envPtr)->DeleteLocalRef(envPtr, inputStreamClass);
        (*envPtr)->DeleteLocalRef(envPtr, byteArray);
        copyToError( "Failed to get byte array address", errorMsg);
        return NULL;
    }
    char* buffer = (char*)malloc(bytesRead * sizeof(char));
    memcpy(buffer, byteBuffer, bytesRead);

    // Release memory
    (*envPtr)->DeleteLocalRef(envPtr, inputStreamClass);
    (*envPtr)->DeleteLocalRef(envPtr, byteArray);
    (*envPtr)->ReleaseByteArrayElements(envPtr, byteArray, byteBuffer, 0);
    // Set the count
    *count = (int)bytesRead;
    return buffer;
}

void writeToOutputStream(uintptr_t env, jobject outputStream, const char* buffer, int size, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;

    // Get the OutputStream class
    jclass outputStreamClass = (*envPtr)->GetObjectClass(envPtr, outputStream);
    if (outputStreamClass == NULL) {
        copyToError( "Failed to get the OutputStream class", errorMsg);
        return;
    }

    // Get the write method from OutputStream
    jmethodID writeMethod = (*envPtr)->GetMethodID(envPtr, outputStreamClass, "write", "([B)V");
    if (writeMethod == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, outputStreamClass);
        copyToError("Failed to get write method", errorMsg);
        return;
    }

    // Create a byte array from the buffer
    jbyteArray byteArray = (*envPtr)->NewByteArray(envPtr, size);
    if (byteArray == NULL) {
        (*envPtr)->DeleteLocalRef(envPtr, outputStreamClass);
        copyToError("Failed to create byte array", errorMsg);
        return;
    }
    (*envPtr)->SetByteArrayRegion(envPtr, byteArray, 0, size, (jbyte*)buffer);

    // Call the write method
    (*envPtr)->CallVoidMethod(envPtr, outputStream, writeMethod, byteArray);

    // Release memory
    (*envPtr)->DeleteLocalRef(envPtr, outputStreamClass);
    (*envPtr)->DeleteLocalRef(envPtr, byteArray);
    return;
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

bool enableBluetooth(uintptr_t env, jobject bluetoothAdapter, char** errorMsg) {
    JNIEnv* envPtr = (JNIEnv*)env;
    // Get the enable method
    jclass bluetoothAdapterClass = (*envPtr)->GetObjectClass(envPtr, bluetoothAdapter);
    jmethodID enableMethod = (*envPtr)->GetMethodID(envPtr, bluetoothAdapterClass, "enable", "()Z");
    if (enableMethod == NULL) {
        copyToError("Failed to find enable method", errorMsg);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        return false;
    }
    // Invoke the enable method
    jboolean result = (*envPtr)->CallBooleanMethod(envPtr, bluetoothAdapter, enableMethod);
    if ((*envPtr)->ExceptionCheck(envPtr)) {
        copyToError("Failed to invoke enable method", errorMsg);
        (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
        return false;
    }
    // Release memory
    (*envPtr)->DeleteLocalRef(envPtr, bluetoothAdapterClass);
    return result;
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
	_ "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
	"unsafe"
	_ "unsafe"
)

// adapterAndroid has references android.bluetooth.BluetoothAdapter in java on android
type adapterAndroid struct {
	javaAdapter C.jobject
	stop        chan struct{}
}

// socketAndroid is representing bluetooth socket on android
type socketAndroid struct {
	javaSocket C.jobject
	rw         *readWriterAndroid
}

// serverSocketAndroid is representing bluetooth server socket on android
type serverSocketAndroid struct {
	javaServerSocket C.jobject
}

// readWriterAndroid has java input and output stream on android
type readWriterAndroid struct {
	in, out C.jobject
}

// runOnJVM is link to fyne.io/fyne/v2/internal/driver/mobile/mobileinit.RunOnJVM
// it direct call internal fyne RunOnJVM
func runOnJVM(fn func(vm, env, ctx uintptr) error) error {
	return driver.RunNative(func(context interface{}) error {
		if androidContext, ok := context.(driver.AndroidContext); ok {
			return fn(androidContext.VM, androidContext.Env, androidContext.Ctx)
		}
		return errors.New("no context android")
	})
}

// NewBluetoothDefaultAdapter get Bluetooth adapter
func NewBluetoothDefaultAdapter() (_ Adapter, e error) {
	var bb *adapterAndroid
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		adapter := C.getBluetoothAdapter(C.uintptr_t(env), &errMsgC)
		if errMsgC != nil {
			er := errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			e = er
			return nil
		}
		bb = &adapterAndroid{javaAdapter: adapter, stop: make(chan struct{})}
		return nil
	})
	enable, err := bb.enable()
	if err != nil {
		return nil, errors.Join(err, bb.Close())
	}
	if !enable {
		return nil, errors.Join(errors.New("not available bluetooth"), bb.Close())
	}
	return bb, errors.Join(e, err)
}

func (a *adapterAndroid) FetchAddress() (str string, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		nameC := C.getBluetoothName(C.uintptr_t(env), a.javaAdapter, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		str = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return str, errors.Join(e, err)
}

// Close clean bluetooth adapter from memory
func (a *adapterAndroid) Close() (e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		C.freeBluetoothAdapter(C.uintptr_t(env), a.javaAdapter, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		return nil
	})
	return errors.Join(e, err)
}

// GetBluetoothServerSocket returns ServerSocket which is listening on bluetooth
func (a *adapterAndroid) GetBluetoothServerSocket() (bs ServerSocket, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		socket := C.createBluetoothServer(C.uintptr_t(env), a.javaAdapter, &errMsgC)
		if errMsgC != nil {
			err := errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			bs, e = nil, err
			return nil
		}
		bs = &serverSocketAndroid{javaServerSocket: socket}
		return nil
	})
	return bs, errors.Join(e, err)
}

// ConnectAsClientToServer take MAC address of server and
// return conection in BluetoothSocket
func (a *adapterAndroid) ConnectAsClientToServer(address string) (bs Socket, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		caddress := C.CString(address)
		defer C.free(unsafe.Pointer(caddress))
		sock := C.connectToBluetoothServer(C.uintptr_t(env), a.javaAdapter, caddress, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		bs = &socketAndroid{javaSocket: sock}
		return nil
	})
	return bs, errors.Join(err, e)
}

func (a *adapterAndroid) enable() (b bool, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		b = bool(C.enableBluetooth(C.uintptr_t(env), a.javaAdapter, &errMsgC))
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		return nil
	})
	return b, errors.Join(e, err)
}

// FetchStringData it is usefully if GetAddress return empty string,
// it try to set internal address
func (b *serverSocketAndroid) FetchStringData() (res string, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		nameC := C.getAddress(C.uintptr_t(env), b.javaServerSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		res = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return res, errors.Join(e, err)
}

func (b *serverSocketAndroid) Close() (e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		C.closeBluetoothServerSocket(C.uintptr_t(env), b.javaServerSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		return nil
	})
	return errors.Join(e, err)
}

// Accept accepting client and return conection in BluetoothSocket
func (b *serverSocketAndroid) Accept() (bs Socket, e error) {
	var errMsgC *C.char
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		sock := C.acceptBluetoothClient(C.uintptr_t(env), b.javaServerSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		bs = &socketAndroid{javaSocket: sock}
		return nil
	})
	return bs, errors.Join(err, e)
}

func (b *socketAndroid) FetchStringData() (string, error) {
	name, err := b.fetchName()
	if err != nil {
		return "", err
	}
	addr, err := b.fetchAddress()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("name:%s, address: %s", name, addr), nil
}

func (b *socketAndroid) Read(bytes []byte) (int, error) {
	if b.rw == nil {
		var err error
		b.rw, err = b.getReadWriter()
		if err != nil {
			return 0, err
		}
	}
	return b.rw.Read(bytes)
}

func (b *socketAndroid) Write(bytes []byte) (int, error) {
	if b.rw == nil {
		var err error
		b.rw, err = b.getReadWriter()
		if err != nil {
			return 0, err
		}
	}
	return b.rw.Write(bytes)
}

func (b *socketAndroid) fetchName() (res string, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		nameC := C.getClientName(C.uintptr_t(env), b.javaSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		res = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return res, errors.Join(err, e)
}

func (b *socketAndroid) fetchAddress() (res string, e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		nameC := C.getClientAddress(C.uintptr_t(env), b.javaSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		res = C.GoString(nameC)
		C.free(unsafe.Pointer(nameC))
		return nil
	})
	return res, errors.Join(err, e)
}

// Close is closing socket
func (b *socketAndroid) Close() (e error) {
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		C.closeBluetoothServerSocket(C.uintptr_t(env), b.javaSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		return nil
	})

	return errors.Join(e, err, b.rw.Close())
}

func (b *socketAndroid) getReadWriter() (rw *readWriterAndroid, e error) {
	var errMsgC *C.char
	rw = &readWriterAndroid{}
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		inputStream := C.getBluetoothInputStream(C.uintptr_t(env), b.javaSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			return nil
		}
		outputStream := C.getBluetoothOutputStream(C.uintptr_t(env), b.javaSocket, &errMsgC)
		if errMsgC != nil {
			e = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			_ = rw.close(inputStream)
			return nil
		}
		rw.in = inputStream
		rw.out = outputStream
		return nil
	})
	return rw, errors.Join(err, e)

}

func (r *readWriterAndroid) Read(p []byte) (n int, err error) {
	if p == nil || len(p) == 0 {
		return -1, errors.New("empty buffer")
	}
	er := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		var result C.int
		dataC := C.readFromInputStream(C.uintptr_t(env), r.in, C.int(cap(p)), &result, &errMsgC)
		if errMsgC != nil {
			err = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
			n = -1
			return nil
		}
		n = int(result)
		if n < 1 {
			return nil
		}
		dataGo := C.GoBytes(unsafe.Pointer(dataC), C.int(n))
		copy(p, dataGo)
		C.free(unsafe.Pointer(dataC))
		return nil
	})
	err = errors.Join(er, err)
	return
}

func (r *readWriterAndroid) Write(p []byte) (n int, err error) {
	if p == nil || len(p) == 0 {
		return 0, errors.New("empty buffer")
	}
	er := runOnJVM(func(vm, env, ctx uintptr) error {
		var errMsgC *C.char
		C.writeToOutputStream(C.uintptr_t(env), r.in, (*C.char)(unsafe.Pointer(&p[0])), C.int(cap(p)), &errMsgC)
		if errMsgC != nil {
			err = errors.New(C.GoString(errMsgC))
			C.free(unsafe.Pointer(errMsgC))
		}
		n = len(p)
		return nil
	})
	return n, errors.Join(er, err)
}

func (r *readWriterAndroid) Close() error {
	return errors.Join(r.close(r.in), r.close(r.out))
}

func (r *readWriterAndroid) close(stream C.jobject) (e error) {
	var errMsgC *C.char
	err := runOnJVM(func(vm, env, ctx uintptr) error {
		C.closeStream0(C.uintptr_t(env), stream, &errMsgC)
		return nil
	})
	if errMsgC != nil {
		e = errors.New(C.GoString(errMsgC))
		C.free(unsafe.Pointer(errMsgC))
	}
	return errors.Join(err, e)
}
