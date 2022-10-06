package main

type Node struct {
	val  interface{}
	next *Node
	prev *Node
}

type ListType struct {
	EqualFunc func(a, b interface{}) bool
}

type List struct {
	ListType
	head *Node
	tail *Node
}

func ListCreate(listType ListType) *List {
	var list List
	list.ListType = listType
	return &list
}

func (list *List) Append(val interface{}) {
	var n Node
	n.val = val
	if list.head == nil {
		list.head = &n
		list.tail = &n
	} else {
		n.prev = list.tail
		list.tail.next = &n
		list.tail = list.tail.next
	}
}

func (list *List) Remove(val interface{}) {
	p := list.head
	for p != nil {
		if list.EqualFunc(p.val, val) {
			break
		}
		p = p.next
	}
	if p != nil {
		p.prev = p.next
		if p.next != nil {
			p.next.prev = p.prev
		}
		p.prev = nil
		p.next = nil
	}
}
