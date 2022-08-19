package main
import "C"
import (
	"github.com/OperatorFoundation/shapeshifter-transports/transports/obfs4/v2"
	"github.com/outspace/obfs4-plug-transport/base"
	"unsafe"
	"net"
	"reflect"
)

var obfs4Clients = map[int]*obfs4.Transport{}
var obfs4Connections = map[int]net.Conn{}
var nextID = 0
var obfs4_transport_listener base.TransportListener = nil

//export Initialize_obfs4_c_client
func Initialize_obfs4_c_client(certString *C.char, iatMode int) (clientKey int) {	
	var err error
	goCertString := C.GoString(certString)
	obfs4Client, err := obfs4.NewObfs4Client(goCertString, iatMode, nil)
	obfs4Clients[nextID] = obfs4Client
	if err != nil {
		return -1
	}
	// This is the return value
	clientKey = nextID

	nextID += 1
	return
}

//export Obfs4_listen
func Obfs4_listen(address_string *C.char) {

	//goAddressString := C.GoString(address_string)
	//obfs4_transport_listener = Obfs4_c_client.Listen(goAddressString)
}

//export Obfs4_dial
func Obfs4_dial(client_id int, address_string *C.char) int {

	goAddressString := C.GoString(address_string)
	var err error
	var transport = obfs4Clients[client_id]
	obfs4_transport_connection, err := transport.Dial(goAddressString)
	if err != nil {
		return -1
	}
	
	if obfs4_transport_connection == nil {
		return 1
	} else {
		obfs4Connections[client_id] = obfs4_transport_connection
		return 0
	}
}

//export Obfs4_write
func Obfs4_write(client_id int, buffer unsafe.Pointer, buffer_length C.int) int {
	var connection = obfs4Connections[client_id]
	var bytesBuffer = C.GoBytes(buffer, buffer_length)
	numberOfBytesWritten, error := connection.Write(bytesBuffer)

	if error != nil {
		return -1
	} else {
		return numberOfBytesWritten
	}
}

//export Obfs4_read
func Obfs4_read(client_id int, buffer unsafe.Pointer, buffer_length int) int {

	var connection = obfs4Connections[client_id]
	if connection == nil {
		return -1
	}
	header := reflect.SliceHeader{uintptr(buffer), buffer_length, buffer_length}
	bytesBuffer := *(*[]byte)(unsafe.Pointer(&header))

	numberOfBytesRead, error := connection.Read(bytesBuffer)

	if error != nil {
		return -1
	} else {
		return numberOfBytesRead
	}
}

//export Obfs4_close_connection
func Obfs4_close_connection(client_id int) {

	var connection = obfs4Connections[client_id]
	connection.Close()
	delete(obfs4Connections, client_id)
	delete(obfs4Clients, client_id)
}

func main() {}
