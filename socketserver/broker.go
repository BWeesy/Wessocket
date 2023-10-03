package main

import (
	"log"
)

type broker struct {
	stopCh    chan struct{}
	publishCh chan message
	subCh     chan chan message
	unsubCh   chan chan message
}

func newBroker() *broker {
	return &broker{
		stopCh:    make(chan struct{}),
		publishCh: make(chan message, 1),
		subCh:     make(chan chan message, 1),
		unsubCh:   make(chan chan message, 1),
	}
}

func (b *broker) start() {
	subscriptions := map[chan message]struct{}{}
	log.Println("Broker running")
	for {
		select {
		case <-b.stopCh:
			return
		case msgCh := <-b.subCh:
			subscriptions[msgCh] = struct{}{}
			log.Println("Subbing - Subscriptions: ", len(subscriptions))
		case msgCh := <-b.unsubCh:
			delete(subscriptions, msgCh)
			log.Println("Unsubbing - Subscriptions: ", len(subscriptions))
		case msg := <-b.publishCh:
			for msgCh := range subscriptions {
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *broker) stop() {
	close(b.stopCh)
}

func (b *broker) subscribe() chan message {
	msgCh := make(chan message, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *broker) unsubscribe(msgCh chan message) {
	b.unsubCh <- msgCh
}

func (b *broker) publish(msg message) {
	b.publishCh <- msg
}
