package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// func main() {

// 	http.HandleFunc("/upscale", FileProcessor)

// 	// Start the HTTP server!
// 	fmt.Println("HTTP server listening on 42070")
// 	if err := http.ListenAndServe("0.0.0.0:42070", nil); err != nil {
// 		fmt.Println(err.Error())
// 	}
// }

type functionHelper struct {
	db *sql.DB
}

func main() {
	// Define command-line flags
	port := flag.Int("port", 8080, "the port number to run the server on")
	dbfile := flag.String("db", "./scalar.db", "SQLite database location")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	fh := &functionHelper{
		db: db,
	}

	endpoints := fh.Endpoints()

	for _, endpoint := range endpoints {
		http.HandleFunc(endpoint.Pattern, endpoint.Handler)
	}

	addr := ":" + strconv.Itoa(*port)
	log.Printf("Server listening on %s", addr)
	fmt.Println(http.ListenAndServe(addr, nil))
}
