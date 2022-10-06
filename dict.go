package main

type entry struct {
	key  interface{}
	val  interface{}
	next *entry
}

type htable struct {
	table []*entry
	size  int64
	mask  int64
	used  int64
}

type DictType struct {
	HashFunc  func(key interface{}) int
	EqualFunc func(k1, k2 interface{}) bool
}

type Dict struct {
	DictType
	HTable    [2]htable
	rehashidx int
	// iterators
}

func DictCreate(dictType DictType) *Dict {
	var dict Dict
	dict.DictType = dictType
	return &dict
}

func (dict *Dict) RandomGet() (key, val interface{}) {
	//TODO:
	return nil, nil
}

func (dic *Dict) RemoveKey(key interface{}) {
	//TODO:
}
