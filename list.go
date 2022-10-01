package main

type Node struct {
	val  *Gobj
	next *Node
	prev *Node
}

type List struct {
	head *Node
	tail *Node
	// todo
}
