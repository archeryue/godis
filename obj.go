package main

type Gtype uint8

const (
	GSTR  Gtype = 0x01
	GLIST Gtype = 0x02
	GDICT Gtype = 0x03
)

type Gval interface{}

type Gobj struct {
	type_    Gtype
	val_     Gval
	refCount int
}
