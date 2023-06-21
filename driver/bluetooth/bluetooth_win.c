#include <winsock2.h>
#include <ws2bth.h>
#include <BluetoothAPIs.h>

// Vytvorenie Bluetooth soketu pre server
SOCKET serverSock = socket(AF_BTH, SOCK_STREAM, BTHPROTO_RFCOMM);

// Adresa a port, na ktorom bude server počúvať
BTH_ADDR localAddress = 0;  // Lokálna adresa Bluetooth zariadenia
int localPort = ...;  // Nastavte správny port

// Viazanie soketu na lokálnu adresu a port
SOCKADDR_BTH serverAddr;
memset(&serverAddr, 0, sizeof(serverAddr));
serverAddr.addressFamily = AF_BTH;
serverAddr.btAddr = localAddress;
serverAddr.port = localPort;

int result = bind(serverSock, (SOCKADDR*)&serverAddr, sizeof(serverAddr));
if (result == SOCKET_ERROR) {
    // Chyba pri viazaní soketu na adresu a port
    int error = WSAGetLastError();
    // Spracovanie chyby
}

// Počúvanie na sokete a čakanie na pripojenie klienta
result = listen(serverSock, SOMAXCONN);
if (result == SOCKET_ERROR) {
    // Chyba pri počúvaní na sokete
    int error = WSAGetLastError();
    // Spracovanie chyby
}

// Pripravený na prijatie pripojenia
SOCKADDR_BTH clientAddr;
int clientAddrSize = sizeof(clientAddr);
SOCKET clientSock = accept(serverSock, (SOCKADDR*)&clientAddr, &clientAddrSize);
if (clientSock == INVALID_SOCKET) {
    // Chyba pri prijímaní pripojenia
    int error = WSAGetLastError();
    // Spracovanie chyby
}

// Teraz môžete použiť clientSock na komunikáciu s klientom

// Zatvorenie soketov po skončení
closesocket(clientSock);
closesocket(serverSock);
