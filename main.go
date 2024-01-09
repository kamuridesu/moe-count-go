package main

import (
	"fmt"
	"os"
	"strings"
)

const imageBasePath = "static/images"

var mainDatabase *Database
var loadedImages *[][]byte
var dbFile string

func argparse() {
	dbFile = "./db/users.db"
	for _, arg := range os.Args {
		if strings.Contains(arg, "-dbfile") {
			dbFile = strings.Split(arg, "=")[1]
		}
	}
}

func init() {
	argparse()

	var err error
	fmt.Println(dbFile)
	mainDatabase, err = StartDB("sqlite3", dbFile)
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
