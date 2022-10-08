package main

import (
	"hash/fnv"
	"log"
	"os"
	"time"
)

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
	args  []*Gobj
	reply *List
}

type CommandProc func(c *GodisClient)

// do not support bulk command
type GodisCommand struct {
	name  string
	proc  CommandProc
	arity int
}

// Global Varibles
var server GodisServer
var cmdTable []GodisCommand = []GodisCommand{
	{"get", getCommand, 2},
	{"set", setCommand, 3},
	//TODO
}

func getCommand(c *GodisClient) {
	//TODO
}

func setCommand(c *GodisClient) {
	//TODO
}

func ReadQueryFromClient(loop *AeLoop, fd int, extra interface{}) {
	//TODO: read query from client
	//TODO: handle query -> args
	//TODO: proccess command
}

func ClientEqual(a, b interface{}) bool {
	c1, ok := a.(*GodisClient)
	if !ok {
		return false
	}
	c2, ok := a.(*GodisClient)
	if !ok {
		return false
	}
	return c1.fd == c2.fd
}

func GStrEqual(a, b interface{}) bool {
	o1, ok := a.(*Gobj)
	if !ok || o1.Type_ != GSTR {
		return false
	}
	o2, ok := a.(*Gobj)
	if !ok || o1.Type_ != GSTR {
		return false
	}
	return o1.Val_.(string) == o2.Val_.(string)
}

func GStrHash(key interface{}) int {
	o, ok := key.(*Gobj)
	if !ok || o.Type_ != GSTR {
		return 0
	}
	hash := fnv.New32()
	hash.Write([]byte(o.Val_.(string)))
	return int(hash.Sum32())
}

func CreateClient(fd int) *GodisClient {
	var client GodisClient
	client.fd = fd
	client.db = server.db
	client.reply = ListCreate(ListType{EqualFunc: GStrEqual})
	server.aeLoop.AddFileEvent(fd, AE_READABLE, ReadQueryFromClient, nil)
	return &client
}

func AcceptHandler(loop *AeLoop, fd int, extra interface{}) {
	nfd, err := Accept(fd)
	if err != nil {
		log.Printf("accept err: %v\n", err)
		return
	}
	client := CreateClient(nfd)
	//TODO: check max clients limit
	server.clients.Append(client)
}

const EXPIRE_CHECK_COUNT int = 100

// background job, runs every 100ms
func ServerCron(loop *AeLoop, id int, extra interface{}) {
	for i := 0; i < EXPIRE_CHECK_COUNT; i++ {
		key, val := server.db.expire.RandomGet()
		if key == nil {
			break
		}
		if int64(val.(*Gobj).IntVal()) < time.Now().Unix() {
			server.db.data.RemoveKey(key)
			server.db.expire.RemoveKey(key)
		}
	}
}

func initServer(config *Config) error {
	server.port = config.Port
	server.clients = ListCreate(ListType{EqualFunc: ClientEqual})
	server.db = &GodisDB{
		data:   DictCreate(DictType{HashFunc: GStrHash, EqualFunc: GStrEqual}),
		expire: DictCreate(DictType{HashFunc: GStrHash, EqualFunc: GStrEqual}),
	}
	var err error
	if server.aeLoop, err = AeLoopCreate(); err != nil {
		return err
	}
	server.fd, err = TcpServer(server.port)
	return err
}

func main() {
	path := os.Args[1]
	config, err := LoadConfig(path)
	if err != nil {
		log.Printf("config error: %v\n", err)
	}
	err = initServer(config)
	if err != nil {
		log.Printf("init server error: %v\n", err)
	}
	server.aeLoop.AddFileEvent(server.fd, AE_READABLE, AcceptHandler, nil)
	server.aeLoop.AddTimeEvent(AE_NORMAL, 100, ServerCron, nil)
	log.Println("godis server is up.")
	server.aeLoop.AeMain()
}
