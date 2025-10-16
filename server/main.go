package main

import (
	"fmt"
	"net"
	"os"

	"github.com/vini464/distributed_system/communication"
)

var BROKER = os.Getenv("BROKER")
func main() {
	fmt.Println("Server is starting")
	net.Dial(communication.SERVERTYPE, net.JoinHostPort(BROKER, communication.BROKERPORT))
}
