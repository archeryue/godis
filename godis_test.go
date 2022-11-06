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
