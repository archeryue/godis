package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDict(t *testing.T) {
	dict := DictCreate(DictType{HashFunc: GStrHash, EqualFunc: GStrEqual})
	entry := dict.RandomGet()
	assert.Nil(t, entry)

	k1 := CreateObject(GSTR, "k1")
	v1 := CreateObject(GSTR, "v1")
	e := dict.Add(k1, v1)
	assert.Nil(t, e)
	k1.IncrRefCount()
	v1.IncrRefCount()

	entry = dict.Find(k1)
	assert.Equal(t, k1, entry.Key)
	assert.Equal(t, v1, entry.Val)

	e = dict.Delete(k1)
	assert.Nil(t, e)
	entry = dict.Find(k1)
	assert.Nil(t, entry)
}

func TestRehash(t *testing.T) {
	dict := DictCreate(DictType{HashFunc: GStrHash, EqualFunc: GStrEqual})
	entry := dict.RandomGet()
	assert.Nil(t, entry)

	valve := int(INIT_SIZE * (FORCE_RATION + 1))
	for i := 0; i < valve; i++ {
		key := CreateObject(GSTR, fmt.Sprintf("k%v", i))
		val := CreateObject(GSTR, fmt.Sprintf("v%v", i))
		e := dict.Add(key, val)
		assert.Nil(t, e)
	}
	assert.Equal(t, false, dict.isRehashing())

	key := CreateObject(GSTR, fmt.Sprintf("k%v", valve))
	val := CreateObject(GSTR, fmt.Sprintf("v%v", valve))
	e := dict.Add(key, val)
	assert.Nil(t, e)
	assert.Equal(t, true, dict.isRehashing())
	assert.Equal(t, int64(0), dict.rehashidx)
	assert.Equal(t, int64(INIT_SIZE), dict.hts[0].size)
	assert.Equal(t, int64(INIT_SIZE*GROW_RATION), dict.hts[1].size)

	for i := 0; i <= int(INIT_SIZE); i++ {
		dict.RandomGet()
	}
	assert.Equal(t, false, dict.isRehashing())
	assert.Equal(t, int64(INIT_SIZE*GROW_RATION), dict.hts[0].size)
	assert.Nil(t, dict.hts[1])
	for i := 0; i <= valve; i++ {
		key := CreateObject(GSTR, fmt.Sprintf("k%v", i))
		entry = dict.Find(key)
		assert.NotNil(t, entry)
		assert.Equal(t, fmt.Sprintf("v%v", i), entry.Val.StrVal())
	}
}
