
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