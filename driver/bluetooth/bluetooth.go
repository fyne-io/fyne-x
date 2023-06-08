package bluetooth

type Adapter interface {
	FetchAddress() (string, error)
	Close() error
	GetBluetoothServerSocket() (ServerSocket, error)
	ConnectAsClientToServer(string) (Socket, error)
}

type Socket interface {
	FetchStringData() (string, error)
	Close() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type ServerSocket interface {
	FetchStringData() (string, error)
	Close() error
	Accept() (Socket, error)
}
