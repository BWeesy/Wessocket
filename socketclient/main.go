package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	connection := createConnection(scanner)
	defer connection.Close()

	log.Println("Enter text to send to other users:")

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go handleWrites(connection, scanner, wg)
	go handleReads(connection, wg)

	wg.Wait()
}

func createConnection(scanner *bufio.Scanner) net.Conn {
	config := getConfig()
	log.Printf("Got config %s", config)

	log.Println("Start by declaring your name:")
	scanner.Scan()
	name := scanner.Text()
	connection, err := net.Dial(config.server_type, config.getServerAddress())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection created:")
	go sendMessage(connection, name)
	return connection
}

func handleReads(connection net.Conn, wg *sync.WaitGroup) {
	buffer := make([]byte, 1024)
	defer wg.Done()
	for {
		mLen, readErr := connection.Read(buffer)

		if readErr != nil {
			log.Println("Error reading: ", readErr.Error())
			break
		}
		response := string(buffer[:mLen])

		log.Printf(response)
	}
}

func handleWrites(connection net.Conn, scanner *bufio.Scanner, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for scanner.Scan() {
			text := scanner.Text()
			go sendMessage(connection, text)
		}
		if err := scanner.Err(); err != nil {
			log.Println("Error writing:", err.Error())
		}
	}
}

func sendMessage(connection net.Conn, text string) {
	message := []byte(text)
	writeMessage(connection, message)
}

func writeMessage(connection net.Conn, message []byte) {
	_, writeErr := connection.Write(message)

	if writeErr != nil {
		log.Println("Error writing: ", writeErr.Error())
	}
}
