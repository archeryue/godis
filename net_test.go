package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func EchoServer(s, c, e chan struct{}) {
	sfd, err := TcpServer(6666)
	if err != nil {
		fmt.Printf("tcp server error: %v\n", err)
	}
	fmt.Println("server started")
	s <- struct{}{}
	<-c
	cfd, err := Accept(sfd)
	fmt.Printf("accepted cfd: %v\n", cfd)
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
	e <- struct{}{}
}

//func TestNet(t *testing.T) {
func useTestAeInstead(t *testing.T) {
	fmt.Println("test net lib, rerun the test if the program is blocked")
	s := make(chan struct{})
	c := make(chan struct{})
	e := make(chan struct{})
	go EchoServer(s, c, e)
	<-s
	host := [4]byte{127, 0, 0, 1}
	cfd, err := Connect(host, 6666)
	fmt.Printf("connected cfd: %v\n", cfd)
	time.Sleep(100 * time.Millisecond)
	c <- struct{}{}
	assert.Nil(t, err)
	msg := "helloworld"
	n, err := Write(cfd, []byte(msg))
	assert.Nil(t, err)
	assert.Equal(t, 10, n)
	<-e
	buf := make([]byte, 10)
	n, err = Read(cfd, buf)
	assert.Nil(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, msg, string(buf))
}
