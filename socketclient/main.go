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

	// Read input line by line
	for scanner.Scan() {
		text := scanner.Text() // Get the current line of text
		if text == "" {
			break // Exit loop if an empty line is entered
		}
		go sendMessage(text)
	}

	log.Println("Entry finished, thanks")

	if err := scanner.Err(); err != nil {
		log.Println("Error:", err)
	}
}

func sendMessage(text string) {
	connection, err := net.Dial(SERVER_TYPE, SERVER_ADDRESS)
	if err != nil {
		panic(err)
	}

	message := []byte(text)
	log.Println("Sending: " + text)

	writeMessage(connection, message)
	readResponse(connection)

	defer connection.Close()
}

func readResponse(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, readErr := connection.Read(buffer)

	if readErr != nil {
		log.Println("Error reading: ", readErr.Error())
	}

	log.Printf("Recieved %v Bytes with message: %v", mLen, string(buffer[:mLen]))
}

func writeMessage(connection net.Conn, message []byte) {
	_, writeErr := connection.Write(message)

	if writeErr != nil {
		log.Println("Error writing: ", writeErr.Error())
	}
}
