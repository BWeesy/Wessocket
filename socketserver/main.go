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

func main() {
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
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		log.Println("Error reading: ", err.Error())
	}

	log.Println("Recieved: ", string(buffer[:mLen]))
	_, err = connection.Write([]byte("Ack, recieved:" + string(buffer[:mLen])))
	connection.Close()
}
