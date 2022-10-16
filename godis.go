package main

import (
	"errors"
	"hash/fnv"
	"log"
	"os"
	"time"
)

type CmdType = byte

const (
	COMMAND_UNKNOWN CmdType = 0x00
	COMMAND_INLINE  CmdType = 0x01
	COMMAND_BULK    CmdType = 0x02
)

const (
	GODIS_IO_BUF   int = 1024 * 16
	GODIS_MAX_BULK int = 1024 * 4
)

type GodisDB struct {
	data   *Dict
	expire *Dict
}

type GodisServer struct {
	fd      int
	port    int
	db      *GodisDB
	clients map[int]*GodisClient
	aeLoop  *AeLoop
}

type GodisClient struct {
	fd       int
	db       *GodisDB
	args     []*Gobj
	reply    *List
	queryBuf []byte
	queryLen int
	cmdTy    CmdType
	bulkNum  int
	bulkLen  int
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

func processCommand(c *GodisClient) {
	//TODO: lookup command
	//TODO: call command
	//TODO: decrRef args
}

func freeClient(client *GodisClient) {
	//TODO: delete file event
	//TODO: decrRef reply & args list
	//TODO: delete from clients
}

func resetClient(client *GodisClient) {

}

func handleInlineBuf(client *GodisClient) (bool, error) {
	return false, nil
}

func handleBulkBuf(client *GodisClient) (bool, error) {
	return false, nil
}

func handleQueryBuf(client *GodisClient) error {
	for client.queryLen > 0 {
		if client.cmdTy == COMMAND_UNKNOWN {
			if client.queryBuf[0] == '*' {
				client.cmdTy = COMMAND_BULK
			} else {
				client.cmdTy = COMMAND_INLINE
			}
		}
		// trans query -> args
		var ok bool
		var err error
		if client.cmdTy == COMMAND_INLINE {
			ok, err = handleInlineBuf(client)
		} else if client.cmdTy == COMMAND_BULK {
			ok, err = handleBulkBuf(client)
		} else {
			return errors.New("unknow Godis Command Type")
		}
		if err != nil {
			return err
		}
		// after query -> args
		if ok {
			if len(client.args) == 0 {
				resetClient(client)
			} else {
				processCommand(client)
			}
		} else {
			break
		}
	}
	return nil
}

func ReadQueryFromClient(loop *AeLoop, fd int, extra interface{}) {
	client := extra.(*GodisClient)
	if len(client.queryBuf)-client.queryLen < GODIS_MAX_BULK {
		client.queryBuf = append(client.queryBuf, make([]byte, GODIS_MAX_BULK, GODIS_MAX_BULK)...)
	}
	n, err := Read(fd, client.queryBuf[client.queryLen:])
	if err != nil {
		log.Printf("client %v read err: %v\n", fd, err)
		return
	}
	client.queryLen += n
	err = handleQueryBuf(client)
	if err != nil {
		log.Printf("handle query buf err: %v\n", err)
		return
	}
}

func GStrEqual(a, b *Gobj) bool {
	if a.Type_ != GSTR || b.Type_ != GSTR {
		return false
	}
	return a.Val_.(string) == a.Val_.(string)
}

func GStrHash(key *Gobj) int {
	if key.Type_ != GSTR {
		return 0
	}
	hash := fnv.New32()
	hash.Write([]byte(key.Val_.(string)))
	return int(hash.Sum32())
}

func CreateClient(fd int) *GodisClient {
	var client GodisClient
	client.fd = fd
	client.db = server.db
	client.queryBuf = make([]byte, GODIS_IO_BUF, GODIS_IO_BUF)
	client.reply = ListCreate(ListType{EqualFunc: GStrEqual})
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
	server.clients[fd] = client
	server.aeLoop.AddFileEvent(fd, AE_READABLE, ReadQueryFromClient, client)
}

const EXPIRE_CHECK_COUNT int = 100

// background job, runs every 100ms
func ServerCron(loop *AeLoop, id int, extra interface{}) {
	for i := 0; i < EXPIRE_CHECK_COUNT; i++ {
		key, val := server.db.expire.RandomGet()
		if key == nil {
			break
		}
		if int64(val.IntVal()) < time.Now().Unix() {
			server.db.data.RemoveKey(key)
			server.db.expire.RemoveKey(key)
		}
	}
}

func initServer(config *Config) error {
	server.port = config.Port
	server.clients = make(map[int]*GodisClient)
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
