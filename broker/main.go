package main

import (
	"fmt"
	"net"
)

const (
	SERVERTYPE = "tcp"
	HOSTNAME   = "broker:7575"
)

func main() {
	listener, err := net.Listen("tcp", "broker")
	if err != nil {
		panic(err)
	}

	broker := NewBroker()

	for {
		conn, err := listener.Accept()
		if err != nil {
			go handle_connection(conn, broker)
		}
	}

}

// Essa função lida com a comunicação dos clientes
// [4 bytes - tamanho da mesagem][X bytes - messagem ]
// Client -> Broker
// Subscribe:
// [CMD: Subscribe; TOPIC: topic_name]
// Publish:
// [CMD: Publish; TOPIC: topic_name; MESSAGE: message_body]
//
// Broker -> Client
// [TYPE: status; VALUE: OK/ERROR]
// [TYPE: msg; VALUE: msg_body]
func handle_connection(conn net.Conn, broker *Broker) {
	defer conn.Close()
	sendding_chan := make(chan string)
	receiving_chan := make(chan string)

	// receiving a client message
	go func() {
		for {

		}
	}()

	for {
		select {
		case msg := <-sendding_chan:
			//TODO: send message to client
			fmt.Println(msg)
		case <-broker.Quit:
			fmt.Println("[DEBUG] - Broker is stopped")
			return
		}

	}
}
