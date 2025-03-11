package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/tursodatabase/go-libsql"
)

var (
	db   *sql.DB
	lock = &sync.Mutex{}
)

func ConnectToDb() {
	dbName := "file:./database/local.db"
	log.Println("Connecting to db at " + dbName + "...")
	newDb, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	db = newDb
	log.Println("Connected to db")
}

func Db() *sql.DB {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if db == nil {
			ConnectToDb()
		}
		// else {
		// 	log.Println("Reusing db connection")
		// }
	}
	// else {
	// 	log.Println("Reusing db connection")
	// }
	return db
}
