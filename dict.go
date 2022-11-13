package main

import (
	"errors"
	"math"
)

const INIT_SIZE int64 = 8
const FORCE_RATION int64 = 3
const GROW_RATION int64 = 2

var (
	EP_ERR = errors.New("expand error")
)

type entry struct {
	key  *Gobj
	val  *Gobj
	next *entry
}

type htable struct {
	table []*entry
	size  int64
	mask  int64
	used  int64
}

type DictType struct {
	HashFunc  func(key *Gobj) int64
	EqualFunc func(k1, k2 *Gobj) bool
}

type Dict struct {
	DictType
	hts       [2]*htable
	rehashidx int64
	// iterators
}

func DictCreate(dictType DictType) *Dict {
	var dict Dict
	dict.DictType = dictType
	return &dict
}

func (dict *Dict) isRehashing() bool {
	return dict.rehashidx != -1
}

func (dict *Dict) rehash(step int) {
	for step > 0 {
		if dict.hts[0].used == 0 {
			dict.hts[0] = dict.hts[1]
			dict.hts[1] = nil
			dict.rehashidx = -1
			return
		}
		// find an nonull slot
		for dict.hts[0].table[dict.rehashidx] == nil {
			dict.rehashidx += 1
		}
		// migrate all keys in this slot
		entry := dict.hts[0].table[dict.rehashidx]
		for entry != nil {
			ne := entry.next
			idx := dict.HashFunc(entry.key) & dict.hts[1].mask
			entry.next = dict.hts[1].table[idx]
			dict.hts[1].table[idx] = entry
			dict.hts[0].used -= 1
			dict.hts[1].used += 1
			entry = ne
		}
		dict.hts[0].table[dict.rehashidx] = nil
		dict.rehashidx += 1
		step -= 1
	}
}

func nextPower(size int64) int64 {
	for i := INIT_SIZE; i < math.MaxInt64; i *= 2 {
		if i >= size {
			return i
		}
	}
	return -1
}

func (dict *Dict) expand(size int64) error {
	sz := nextPower(size)
	if dict.isRehashing() || dict.hts[0].used > sz {
		return EP_ERR
	}
	var ht htable
	ht.size = sz
	ht.mask = sz - 1
	ht.table = make([]*entry, sz)
	ht.used = 0
	// check for init
	if dict.hts[0] == nil {
		dict.hts[0] = &ht
		return nil
	}
	// start rehashing
	dict.hts[1] = &ht
	dict.rehashidx = 0
	return nil
}

func (dict *Dict) expandIfNeeded() error {
	if dict.isRehashing() {
		return nil
	}
	if dict.hts[0].size == 0 {
		return dict.expand(INIT_SIZE)
	}
	if (dict.hts[0].used > dict.hts[0].size) && (dict.hts[0].used/dict.hts[0].size > FORCE_RATION) {
		return dict.expand(dict.hts[0].size * GROW_RATION)
	}
	return nil
}

func (dict *Dict) RandomGet() (key, val *Gobj) {
	//TODO:
	return nil, nil
}

func (dic *Dict) RemoveKey(key *Gobj) {
	//TODO:
}
