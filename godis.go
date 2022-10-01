package main

import "fmt"

type GodisDB struct {
	data   *Dict
	expire *Dict
}

type GodisServer struct {
	fd      int
	port    int
	db      *GodisDB
	clients *List
	aeLoop  *AeLoop
}

type GodisClient struct {
	fd    int
	db    *GodisDB
	query string
	reply *List
}

func main() {
	fmt.Println("vim-go")
}
