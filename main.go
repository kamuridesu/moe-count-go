package main

import (
	"database/sql"
	// "fmt"
)

const BASEPATH = "static/images"

var DB *sql.DB
var IMAGES *[][]byte

func init() {
	var err error
	DB, err = StartDB("./db/users.db")
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
