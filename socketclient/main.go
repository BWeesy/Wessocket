package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const (
	SERVER_HOST    = "localhost"
	SERVER_PORT    = "9988"
	SERVER_TYPE    = "tcp"
	SERVER_ADDRESS = SERVER_HOST + ":" + SERVER_PORT
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	log.Println("Enter some text (press Ctrl+D or Ctrl+Z to end):")

	for scanner.Scan() {
		text := scanner.Text()
		for i := 0; i < 100; i++ { // Simulate large traffic spike
			go sendMessage(text)
		}
	}

	log.Println("Entry finished, thanks")

	if err := scanner.Err(); err != nil {
		log.Println("Error:", err)
	}
}

func sendMessage(text string) {
	connection, err := net.Dial(SERVER_TYPE, SERVER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte(text)
	writeMessage(connection, message)
	response := readResponse(connection)

	log.Printf("Sent: %v. Recieved: %v.", text, response)

	defer connection.Close()
}

func readResponse(connection net.Conn) string {
	buffer := make([]byte, 1024)
	mLen, readErr := connection.Read(buffer)

	if readErr != nil {
		log.Println("Error reading: ", readErr.Error())
	}

	return string(buffer[:mLen])
}

func writeMessage(connection net.Conn, message []byte) {
	_, writeErr := connection.Write(message)

	if writeErr != nil {
		log.Println("Error writing: ", writeErr.Error())
	}
}
