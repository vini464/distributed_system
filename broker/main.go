package main

import (
	"fmt"
	"net"

	"github.com/vini464/distributed_system/communication"
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
	communication_chan := make(chan communication.Message)

	// receiving a client message
	go func() {
		var msg communication.Message
		for {
			err := communication.ReceiveMessage(conn, &msg)
			if err != nil {
				return
			}
			communication_chan <- msg
		}
	}()

	for {
		select {
		case msg := <-communication_chan:
			switch msg.Cmd {
			case communication.SUBSCRIBE:
				broker.Subscribe(msg.Tpc, conn)
			case communication.UNSUB:
				broker.Unsubscribe(msg.Tpc, conn)
			case communication.PUBLISH:
				broker.Publish(msg.Tpc, msg.Msg)
			}
		case <-broker.Quit:
			fmt.Println("[DEBUG] - Broker is stopped")
			return
		}
	}
}
