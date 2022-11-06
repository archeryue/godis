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

	ReadQuery(client, "*3\r\n$3\r\nset\r\n$3")
	ok, err = handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, false, ok)

	ReadQuery(client, "\r\nkey\r\n$3\r\nval\r\n")
	ok, err = handleBulkBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, 3, len(client.args))
}

func TestProcessQueryBuf(t *testing.T) {
	client := CreateClient(0)
	ReadQuery(client, "*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	err := ProcessQueryBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(client.args))

	ReadQuery(client, "set key val\r\n")
	err = ProcessQueryBuf(client)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(client.args))
}
