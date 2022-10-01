package main

import (
	"fmt"
	"log"
	"os"
)

type GodisDB struct {
	data   *Dict
	expire *Dict
}

type GodisServer struct {
	fd      int
	port    int
	db      *GodisDB
	clients *List
	aeLoop  *AeLoop
}

type GodisClient struct {
	fd    int
	db    *GodisDB
	query string
	reply *List
}

type CommandProc func(c *GodisClient)

// do not support bulk command
type GodisCommand struct {
	name  string
	proc  CommandProc
	arity int
}

// Global Varibles
var server GodisServer
var cmdTable []GodisCommand

func initServer(config *Config) error {
	//TODO
	return nil
}

func getCommand(c *GodisClient) {
	//TODO
}

func setCommand(c *GodisClient) {
	//TODO
}

func initCmdTable() {
	cmdTable = []GodisCommand{
		{"get", getCommand, 2},
		{"set", setCommand, 3},
		//TODO
	}
}

func main() {
	path := os.Args[1]
	config, err := LoadConfig(path)
	if err != nil {
		log.Printf("config error: %v\n", err)
	}
	err = initServer(config)
	if err != nil {
		log.Printf("init server error: %v\n", err)
	}
	initCmdTable()
	log.Println("godis server is up.")
	server.aeLoop.AeMain()
}
