package main

import (
	"log"
	"net"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type configuration struct {
	server_host string
	server_port string
	server_type string
}

func (c configuration) getServerAddress() string {
	return c.server_host + ":" + c.server_port
}

func getConfig() configuration {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")

	if err != nil {
		log.Println("Unable to find config.yaml, using defaults")
		return configuration{server_host: "localhost", server_port: "9988", server_type: "tcp"}
	}

	return configuration{
		server_host: GetOutboundIP().String(),
		server_port: config.String("server_port"),
		server_type: config.String("server_type"),
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
