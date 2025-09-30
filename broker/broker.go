package main

import (
	"net"
	"sync"

	"github.com/vini464/distributed_system/communication"
)

type Broker struct {
	Mutex       sync.Mutex
	Subscribers map[string][]net.Conn
	Quit        chan struct{} // é um canal de sinal. Ele indica se o Broker foi fechado para as treads
	Closed      bool
}

// Cria um novo broker
func NewBroker() *Broker {
	return &Broker{
		Subscribers: make(map[string][]net.Conn),
		Quit:        make(chan struct{}),
	}
}

func (b *Broker) Close() {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.Closed {
		return
	}

	close(b.Quit)
	b.Closed = true

	for _, topic := range b.Subscribers {
		for _, conn := range topic {
			conn.Close() // closes the connection when the broker is closed
		}
	}
}

// se inscreve em um tópico
func (b *Broker) Subscribe(topic string, conn net.Conn) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.Closed {
		return
	}

	b.Subscribers[topic] = append(b.Subscribers[topic], conn)
}

// remove a inscrição de uma conexão em um topico
func (b *Broker) Unsubscribe(topic string, conn net.Conn) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.Closed {
		return
	}

	for id, c := range b.Subscribers[topic] {
		if c == conn {
			b.Subscribers[topic] = append(b.Subscribers[topic][:id], b.Subscribers[topic][id+1:]...) 
			return
		}
	}
}

// envia a mensagem para o canal inscrito
func (b *Broker) Publish(topic string, msg_body []byte) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.Closed {
		return
	}

	msg := communication.Message{
		Cmd:  communication.MESSAGE,
		Msg: msg_body,
	}

	for _, conn := range b.Subscribers[topic] {
		communication.SendMessage(conn, msg)
	}
}

