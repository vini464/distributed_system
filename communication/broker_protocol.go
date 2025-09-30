package communication

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

// protocol
// [4 bytes - tamanho da mesagem][X bytes - messagem ]
// Client -> Broker
// Subscribe:
// [CMD: Subscribe; TOPIC: topic_name]
// Publish:
// [CMD: Publish; TOPIC: topic_name; MESSAGE: message_body]
//
// Broker -> Client
// [TYPE: status; VALUE: OK/ERROR] -> when a client ask for something (like subscribe or post some message)
// [TYPE: msg; VALUE: msg_body]    -> when the broker has to send a broadcast (a publisher post something, all subs will receive this)

type Request struct {
	Cmd string `json:"cmd"`
	Tpc string `json:"tpc"`
	Msg []byte `json:"msg,omitempty"` // description about this field is on the Server-Client protocol file
}

type Response struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}

type Sendable interface {
	Request | Response
}

// some cmd and types
const (
	UNSUB     = "USUB"
	SUBSCRIBE = "SUB"
	PUBLISH   = "PUB"

	STATUS  = "STA"
	MESSAGE = "MSG"
)

func SendMessage[T Sendable](conn net.Conn, message T) error {
	serialized, err := json.Marshal(message)
	if err != nil {
		return err
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(serialized)))
	_, err = conn.Write(header)
	if err != nil {
		return err
	}
	_, err = conn.Write(serialized)
	return err
}

func ReceiveMessage[T Sendable](conn net.Conn, message *T) error {
	// first: receives the message length
	header := make([]byte, 4)
	bytes_received := 0
	for bytes_received < len(header) {
		readed, err := conn.Read(header[bytes_received:])
		if err != nil {
			return err
		}
		bytes_received += readed
	}
	// second: Receives the data
	data := make([]byte, int(binary.BigEndian.Uint32(header)))
	bytes_received = 0
	for bytes_received < len(data) {
		readed, err := conn.Read(data[bytes_received:])
		if err != nil {
			return err
		}
		bytes_received += readed
	}
	// last but not least: unmarshal message data
	err := json.Unmarshal(data, message)

	return err
}
