package main

import (
	"database/sql"
	"flag"
	"fmt"
)

const imageBasePath = "static/images"

var mainDatabase *sql.DB
var IMAGES *[][]byte
var DBFILE string

func argparse() {
	fileNamePtr := flag.String("dbfile", "./db/users.db", "database filelame with path")
	flag.Parse()
	DBFILE = *fileNamePtr
}

func init() {
	argparse()

	var err error
	fmt.Println(DBFILE)
	mainDatabase, err = StartDB(DBFILE)
	if err != nil {
		panic(err)
	}

	IMAGES, err = LoadAllImages(imageBasePath)
	if err != nil {
		panic(err)
	}
}

func main() {
	serve()
	defer mainDatabase.Close()
}
