package main

import "strconv"

type Gtype uint8

const (
	GSTR  Gtype = 0x00
	GLIST Gtype = 0x01
	GSET  Gtype = 0x02
	GZSET Gtype = 0x03
	GDICT Gtype = 0x04
)

type Gval interface{}

type Gobj struct {
	Type_    Gtype
	Val_     Gval
	refCount int
}

func (o *Gobj) IntVal() int64 {
	if o.Type_ != GSTR {
		return 0
	}
	val, _ := strconv.ParseInt(o.Val_.(string), 10, 64)
	return val
}

func (o *Gobj) StrVal() string {
	if o.Type_ != GSTR {
		return ""
	}
	return o.Val_.(string)
}

func CreateFromInt(val int64) *Gobj {
	return &Gobj{
		Type_:    GSTR,
		Val_:     strconv.FormatInt(val, 10),
		refCount: 1,
	}
}

func CreateObject(typ Gtype, ptr interface{}) *Gobj {
	return &Gobj{
		Type_:    typ,
		Val_:     ptr,
		refCount: 1,
	}
}

func (o *Gobj) IncrRefCount() {
	o.refCount++
}

func (o *Gobj) DecrRefCount() {
	o.refCount--
	if o.refCount == 0 {
		// let GC do the work
		o.Val_ = nil
	}
}
