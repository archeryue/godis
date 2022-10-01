package main

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

type DictType interface {
	HashFunc() int
	CompareFunc() int
}

type Dict struct {
	DictType
	HTable    [2]htable
	rehashidx int
	// iterators
}
