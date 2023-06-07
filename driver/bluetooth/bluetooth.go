package bluetooth

type Adapter interface {
	FetchAddress() (string, error)
	Close() error
	GetBluetoothServerSocket() (ServerSocket, error)
	ConnectAsClientToServer(string) (Socket, error)
}

type ReadWriter interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type Socket interface {
	FetchName() (string, error)
	FetchAddress() (string, error)
	Close() error
	GetReadWriter() (ReadWriter, error)
}

type ServerSocket interface {
	FetchAddress() (string, error)
	Close() error
	Accept() (Socket, error)
}
