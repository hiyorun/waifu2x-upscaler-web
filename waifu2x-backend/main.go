package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	beanstalk "github.com/beanstalkd/go-beanstalk"
	"github.com/hiyorun/waifu2x-upscaler-web/pkg/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type functionHelper struct {
	db           *sql.DB
	beanstalk    *beanstalk.Conn
	sharedFolder string
	wsPool       *websocket.Pool
}

func main() {
	// Define command-line flags
	port := flag.Int("port", 8080, "the port number to run the server on")
	dbfile := flag.String("db", "./scalar.db", "SQLite database location")
	beanstalkAddr := flag.String("beanstalk", "127.0.0.1:11300", "The beanstalk server address")
	sharedFolder := flag.String("sharedFolder", "./", "Shared folder location")

	flag.Parse()

	pool := websocket.NewPool()
	go pool.Start()

	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	beanstalk, err := beanstalk.Dial("tcp", *beanstalkAddr)
	if err != nil {
		fmt.Println("Error connecting to Beanstalkd:", err)
		return
	}
	defer beanstalk.Close()

	fh := &functionHelper{
		db:           db,
		beanstalk:    beanstalk,
		sharedFolder: *sharedFolder,
		wsPool:       pool,
	}

	endpoints := fh.Endpoints()

	for _, endpoint := range endpoints {
		http.HandleFunc(endpoint.Pattern, endpoint.Handler)
	}

	// defer fh.webSocket.Close()

	addr := ":" + strconv.Itoa(*port)
	log.Printf("Server listening on %s", addr)
	fmt.Println(http.ListenAndServe(addr, nil))
}
