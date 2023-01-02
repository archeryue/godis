package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ReadQuery(client *GodisClient, query string) {
	for _, v := range []byte(query) {
		client.queryBuf[client.queryLen] = v
		client.queryLen += 1
	}
}

func TestInlineBuf(t *testing.T) {
	client := CreateClient(0)
	ReadQuery(client, "set key val\r\n")
	ok, err := handleInlineBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)

	ReadQuery(client, "set ")
	ok, err = handleInlineBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	ReadQuery(client, "key ")
	ok, err = handleInlineBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	ReadQuery(client, "val\r\n")
	ok, err = handleInlineBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, len(client.args))
}

func TestBulkBuf(t *testing.T) {
	client := CreateClient(0)

	ReadQuery(client, "*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	ok, err := handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, len(client.args))

	ReadQuery(client, "*3\r")
	ok, err = handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	ReadQuery(client, "\n$3\r\nset\r\n$3")
	ok, err = handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	ReadQuery(client, "\r\nkey\r")
	ok, err = handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	ReadQuery(client, "\n$3\r\nval\r\n")
	ok, err = handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, len(client.args))
}

func TestProcessQueryBuf(t *testing.T) {
	var conf Config
	initServer(&conf)
	// just need real fd to support AddReply
	client := CreateClient(server.fd)
	ReadQuery(client, "*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	err := ProcessQueryBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(client.args))
	key := CreateObject(GSTR, "key")
	val := server.db.data.Get(key)
	assert.Equal(t, "val", val.StrVal())

	ReadQuery(client, "set key val2\r\n")
	err = ProcessQueryBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(client.args))
	val2 := server.db.data.Get(key)
	assert.Equal(t, "val2", val2.StrVal())
}
