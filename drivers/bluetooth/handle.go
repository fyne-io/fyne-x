package bluetooth

type Handle func(readWriter *ReadWriterBluetooth, socketInfo *BluetoothSocket)
