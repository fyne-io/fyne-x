#import <CoreBluetooth/CoreBluetooth.h>

@interface MyPeripheral : NSObject <CBPeripheralManagerDelegate>
@property (nonatomic, strong) CBPeripheralManager *peripheralManager;
@property (nonatomic, strong) CBMutableCharacteristic *characteristic;
@end

@implementation MyPeripheral

- (instancetype)init {
    self = [super init];
    if (self) {
        _peripheralManager = [[CBPeripheralManager alloc] initWithDelegate:self queue:nil];
    }
    return self;
}

- (void)peripheralManagerDidUpdateState:(CBPeripheralManager *)peripheral {
    if (peripheral.state == CBManagerStatePoweredOn) {
        CBUUID *serviceUUID = [CBUUID UUIDWithString:@"YOUR_SERVICE_UUID"];
        CBMutableService *service = [[CBMutableService alloc] initWithType:serviceUUID primary:YES];

        CBUUID *characteristicUUID = [CBUUID UUIDWithString:@"YOUR_CHARACTERISTIC_UUID"];
        CBMutableCharacteristic *characteristic = [[CBMutableCharacteristic alloc] initWithType:characteristicUUID properties:CBCharacteristicPropertyRead | CBCharacteristicPropertyWrite value:nil permissions:CBAttributePermissionsReadable | CBAttributePermissionsWriteable];

        service.characteristics = @[characteristic];
        [self.peripheralManager addService:service];

        NSDictionary *advertisementData = @{CBAdvertisementDataServiceUUIDsKey : @[serviceUUID]};
        [self.peripheralManager startAdvertising:advertisementData];
    }
}

- (void)peripheralManager:(CBPeripheralManager *)peripheral didReceiveReadRequest:(CBATTRequest *)request {
    if ([request.characteristic.UUID isEqual:self.characteristic.UUID]) {
        NSData *value = [@"Hello, World!" dataUsingEncoding:NSUTF8StringEncoding];
        request.value = value;
        [self.peripheralManager respondToRequest:request withResult:CBATTErrorSuccess];
    } else {
        [self.peripheralManager respondToRequest:request withResult:CBATTErrorInvalidHandle];
    }
}

- (void)peripheralManager:(CBPeripheralManager *)peripheral didReceiveWriteRequests:(NSArray<CBATTRequest *> *)requests {
    for (CBATTRequest *request in requests) {
        if ([request.characteristic.UUID isEqual:self.characteristic.UUID]) {
            NSString *receivedString = [[NSString alloc] initWithData:request.value encoding:NSUTF8StringEncoding];
            NSLog(@"Received data: %@", receivedString);
        }
        [self.peripheralManager respondToRequest:request withResult:CBATTErrorSuccess];
    }
}

@end


@interface MyCentral : NSObject <CBCentralManagerDelegate, CBPeripheralDelegate>
@property (nonatomic, strong) CBCentralManager *centralManager;
@property (nonatomic, strong) CBPeripheral *peripheral;
@property (nonatomic, strong) CBCharacteristic *characteristic;
@end

@implementation MyCentral

- (instancetype)init {
    self = [super init];
    if (self) {
        _centralManager = [[CBCentralManager alloc] initWithDelegate:self queue:nil];
    }
    return self;
}

- (void)centralManagerDidUpdateState:(CBCentralManager *)central {
    if (central.state == CBManagerStatePoweredOn) {
        [self.centralManager scanForPeripheralsWithServices:nil options:nil];
    }
}

- (void)centralManager:(CBCentralManager *)central didDiscoverPeripheral:(CBPeripheral *)peripheral advertisementData:(NSDictionary<NSString *,id> *)advertisementData RSSI:(NSNumber *)RSSI {
    if ([peripheral.name isEqualToString:@"YOUR_PERIPHERAL_NAME"]) {
        self.peripheral = peripheral;
        [self.centralManager stopScan];
        [self.centralManager connectPeripheral:peripheral options:nil];
    }
}

- (void)centralManager:(CBCentralManager *)central didConnectPeripheral:(CBPeripheral *)peripheral {
    peripheral.delegate = self;
    [peripheral discoverServices:nil];
}

- (void)peripheral:(CBPeripheral *)peripheral didDiscoverServices:(NSError *)error {
    for (CBService *service in peripheral.services) {
        if ([service.UUID isEqual:[CBUUID UUIDWithString:@"YOUR_SERVICE_UUID"]]) {
            [peripheral discoverCharacteristics:nil forService:service];
        }
    }
}

- (void)peripheral:(CBPeripheral *)peripheral didDiscoverCharacteristicsForService:(CBService *)service error:(NSError *)error {
    for (CBCharacteristic *characteristic in service.characteristics) {
        if ([characteristic.UUID isEqual:[CBUUID UUIDWithString:@"YOUR_CHARACTERISTIC_UUID"]]) {
            self.characteristic = characteristic;
            [peripheral readValueForCharacteristic:characteristic];
        }
    }
}

- (void)peripheral:(CBPeripheral *)peripheral didUpdateValueForCharacteristic:(CBCharacteristic *)characteristic error:(NSError *)error {
    if ([characteristic.UUID isEqual:self.characteristic.UUID]) {
        NSString *receivedString = [[NSString alloc] initWithData:characteristic.value encoding:NSUTF8StringEncoding];
        NSLog(@"Received data: %@", receivedString);
    }
}

- (void)writeData {
    if (self.characteristic && self.peripheral) {
        NSData *value = [@"Hello, Server!" dataUsingEncoding:NSUTF8StringEncoding];
        [self.peripheral writeValue:value forCharacteristic:self.characteristic type:CBCharacteristicWriteWithResponse];
    }
}

@end
