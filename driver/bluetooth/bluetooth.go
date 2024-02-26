package bluetooth

type Adapter interface {
	GetAddress() (string, error)
	Close() error
	GetBluetoothServerSocket() (ServerSocket, error)
	ConnectAsClientToServer(string) (Socket, error)
}

type Socket interface {
	Close() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	StringData() (string, error)
}

type ServerSocket interface {
	Close() error
	Accept() (Socket, error)
}
