package main

import "sync"

type Broker struct {
	Mutex       sync.Mutex
	Subscribers map[string][]chan string
	Quit        chan struct{} // é um canal de sinal. Ele indica se o Broker foi fechado para as treads
	Closed      bool
}

// Cria um novo broker
func NewBroker() *Broker {
	return &Broker{
		Subscribers: make(map[string][]chan string),
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
		for _, ch := range topic {
			close(ch) // fechando os canais de todos os inscritos
		}
	}
}

// se inscreve em um tópico
func (b *Broker) Subscribe(topic string, ch chan string) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.Closed {
		return
	}

	b.Subscribers[topic] = append(b.Subscribers[topic], ch)
}

// envia a mensagem para o canal inscrito
func (b *Broker) Publish(topic string, msg string) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if b.Closed {
		return
	}

	for _, ch := range b.Subscribers[topic] {
		ch <- msg
	}
}
