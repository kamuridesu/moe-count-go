package main

import (
	"database/sql"
	"flag"
	"fmt"
)

const imageBasePath = "static/images"

var mainDatabase *sql.DB
var loadedImages *[][]byte
var dbFile string

func argparse() {
	fileNamePtr := flag.String("dbfile", "./db/users.db", "database filelame with path")
	flag.Parse()
	dbFile = *fileNamePtr
}

func init() {
	argparse()

	var err error
	fmt.Println(dbFile)
	mainDatabase, err = StartDB(dbFile)
	if err != nil {
		panic(err)
	}

	loadedImages, err = LoadAllImages(imageBasePath)
	if err != nil {
		panic(err)
	}
}

func main() {
	serve()
	defer mainDatabase.Close()
}
