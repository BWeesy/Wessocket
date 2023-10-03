package main

import (
	"log"

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
		server_host: config.String("server_host"),
		server_port: config.String("server_port"),
		server_type: config.String("server_type"),
	}
}
