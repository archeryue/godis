package main

import (
	"log"

	"golang.org/x/sys/unix"
)

func TcpServer(port int, addr string) (int, error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("init socket err: %v\n", err)
		return -1, nil
	}
	//TODO: bind & listen
	return s, nil
}
