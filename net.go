package main

import (
	"log"

	"golang.org/x/sys/unix"
)

func TcpServer(port int) (int, error) {
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	defer unix.Close(s)
	if err != nil {
		log.Printf("init socket err: %v\n", err)
		return -1, nil
	}
	err = unix.SetsockoptByte(s, unix.SOL_SOCKET, unix.SO_REUSEADDR, 0x01)
	if err != nil {
		log.Printf("set SO_REUSEADDR err: %v\n", err)
		return -1, nil
	}
	var addr unix.SockaddrInet4
	addr.Port = port
	//TODO: set addr.Addr = any
	err = unix.Bind(s, addr)
	if err != nil {
		log.Printf("bind addr err: %v\n", err)
		return -1, nil
	}
	err = unix.Listen(s, 64)
	if err != nil {
		log.Printf("listen socket err: %v\n", err)
		return -1, nil
	}
	return s, nil
}
