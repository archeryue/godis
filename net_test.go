package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func EchoServer(end chan struct{}) {
	sfd, err := TcpServer(6666)
	if err != nil {
		fmt.Printf("tcp server error: %v\n", err)
	}
	end <- struct{}{}
	cfd, err := Accept(sfd)
	if err != nil {
		fmt.Printf("server accpet error: %v\n", err)
	}
	buf := make([]byte, 10)
	n, err := Read(cfd, buf)
	if err != nil {
		fmt.Printf("server read error: %v\n", err)
	}
	fmt.Printf("read %v bytes\n", n)
	n, err = Write(cfd, buf)
	if err != nil {
		fmt.Printf("server write error: %v\n", err)
	}
	fmt.Printf("write %v bytes\n", n)
	<-end
}

func TestNet(t *testing.T) {
	end := make(chan struct{})
	go EchoServer(end)
	<-end
	host := [4]byte{0, 0, 0, 0}
	cfd, err := Connect(host, 6666)
	assert.Nil(t, err)
	msg := "helloworld"
	n, err := Write(cfd, []byte(msg))
	assert.Nil(t, err)
	assert.Equal(t, 10, n)
	buf := make([]byte, 10)
	n, err = Read(cfd, buf)
	assert.Nil(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, msg, string(buf))
	end <- struct{}{}
}
