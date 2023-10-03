package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type message struct {
	senderName string
	content    string
}

func (m message) toString() string {
	return fmt.Sprintf("%v - %v", m.senderName, m.content)
}

func main() {
	broker := newBroker()
	go broker.start()
	defer broker.stop()

	config := getConfig()
	log.Printf("Got config %s", config)

	server, err := net.Listen(config.server_type, config.getServerAddress())

	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}

	defer server.Close()

	log.Println("Listening on " + config.getServerAddress())
	log.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Fatal("Error Accepting: ", err.Error())
		}

		go handleConnection(connection, broker)
		log.Println("Client Connected")
	}

}

func handleConnection(connection net.Conn, broker *broker) {
	messages := broker.subscribe()
	defer broker.unsubscribe(messages)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go readIncomingAndPublish(connection, broker, wg)
	go writeOutgoingMessages(connection, messages, wg)
	wg.Wait()
}

func writeOutgoingMessages(connection net.Conn, messages chan message, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		message := <-messages

		_, err := connection.Write([]byte(message.toString()))
		if err != nil {
			log.Println("Error writing: ", err.Error())
		}
	}
}

func readIncomingAndPublish(connection net.Conn, broker *broker, wg *sync.WaitGroup) {
	defer wg.Done()
	buffer := make([]byte, 1024)
	defer connection.Close()

	name, namingErr := acceptName(connection)
	if namingErr != nil {
		log.Println("Error setting name: ", namingErr.Error())
		return
	}

	for {
		mLen, err := connection.Read(buffer)
		if err != nil {
			log.Println("Client disconnected: ", err.Error())
			break
		}

		content := string(buffer[:mLen])
		msg := message{name, content}
		log.Println("New Message: ", msg.toString())
		broker.publish(message{name, content})
	}
}

func acceptName(connection net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return "no-name", err
	}
	name := string(buffer[:mLen])
	log.Println("New user joined chat: ", name)
	return name, err
}
