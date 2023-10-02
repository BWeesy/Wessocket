package main

import (
	"fmt"
	"log"
	"net"
)

const (
	SERVER_HOST    = "localhost"
	SERVER_PORT    = "9988"
	SERVER_TYPE    = "tcp"
	SERVER_ADDRESS = SERVER_HOST + ":" + SERVER_PORT
)

type message struct {
	senderName string
	content    string
}

func main() {
	messages := make(chan message)

	fmt.Println("Server running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_ADDRESS)

	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}

	defer server.Close()

	log.Println("Listening on " + SERVER_ADDRESS)
	log.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Fatal("Error Accepting: ", err.Error())
		}

		log.Println("Client Connected")
		go processConnectionForReads(connection, messages)
		go processConnectionForWrites(connection, messages)
	}
}

func processConnectionForWrites(connection net.Conn, messages chan message) {
	for {
		message := <-messages
		log.Printf("Processing message from: %v with content: %v", message.senderName, message.content)
		formatted := fmt.Sprintf("%v - %v", message.senderName, message.content)
		log.Println("Writing message to connection ", string(formatted))

		_, err := connection.Write([]byte(formatted))
		if err != nil {
			log.Println("Error writing: ", err.Error())
		}
	}
}

func processConnectionForReads(connection net.Conn, messages chan message) {
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
		messages <- message{name, content}
		log.Println("Recieved: ", message{name, content})
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
