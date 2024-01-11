package main

import (
	"fmt"
	"os"
	"strings"
)

const imageBasePath = "static/images"

var mainDatabase *Database
var loadedImages *[][]byte
var dbParams string

func getCliArgument(arg string) string {
	return strings.Split(arg, "=")[1]
}

func argparse() {
	dbParams = "./db/users.db"
	database := ""
	username := ""
	password := ""
	dbHost := ""
	dbPort := ""
	for _, arg := range os.Args {
		if strings.Contains(arg, "-version") {
			fmt.Println("0.0.1")
			os.Exit(0)
		} else if strings.Contains(arg, "-dbfile") {
			dbParams = getCliArgument(arg)
		} else if strings.Contains(arg, "-dbname") {
			database = getCliArgument(arg)
		} else if strings.Contains(arg, "-dbuser") {
			username = getCliArgument(arg)
		} else if strings.Contains(arg, "-dbpassword") {
			password = getCliArgument(arg)
		} else if strings.Contains(arg, "-dbhost") {
			dbHost = getCliArgument(arg)
		} else if strings.Contains(arg, "-dbport") {
			dbPort = getCliArgument(arg)
		}
	}
	if dbHost != "" {
		dbParams = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, username, password, database)
	}
}

func init() {
	argparse()

	var err error
	if strings.Contains(dbParams, "host") {
		mainDatabase, err = StartDB("postgres", dbParams)
		if err != nil {
			panic(err)
		}
	} else {
		mainDatabase, err = StartDB("sqlite3", dbParams)
		if err != nil {
			panic(err)
		}
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
