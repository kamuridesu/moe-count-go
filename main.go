package main

import (
	"database/sql"
	"flag"
	"fmt"
)

const BASEPATH = "static/images"

var DB *sql.DB
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
	DB, err = StartDB(DBFILE)
	if err != nil {
		panic(err)
	}

	IMAGES, err = LoadAllImages(BASEPATH)
	if err != nil {
		panic(err)
	}
}

func main() {
	serve()
	defer DB.Close()
}
