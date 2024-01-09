package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseType string

type Database struct {
	db     *sql.DB
	dbType DatabaseType
}

const (
	sqlite3  DatabaseType = "sqlite3"
	postgres DatabaseType = "postgres"
)

func StartDB(__type DatabaseType, parameters string) (*Database, error) {
	var db *sql.DB
	var err error
	switch __type {
	case sqlite3:
		db, err = OpenSqliteDB(parameters)
		if err != nil {
			panic(err)
		}
	case postgres:
		db, err = OpenPostgresDB(parameters)
		if err != nil {
			panic(err)
		}
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS "users" (
		"username"  TEXT NOT NULL,
		"count"     INTEGER NOT NULL,
		PRIMARY KEY("username")
	);`

	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}
	return &Database{db: db, dbType: __type}, nil
}

func executeQuery(db *sql.DB, query string, params ...string) error {
	tx, err := db.Begin()

	if err != nil {
		log.Print(err)
		return err
	}

	stmt, err := tx.Prepare(query)

	if err != nil {
		log.Print(err)
		return err
	}

	defer stmt.Close()

	var args []reflect.Value

	for _, param := range params {
		args = append(args, reflect.ValueOf(param))
	}

	execFun := reflect.ValueOf(stmt.Exec)

	result := execFun.Call(args)

	if result[1].Interface() != nil {
		err := result[1].Interface().(error)
		log.Print(err)
		return err
	}

	err = tx.Commit()

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (db *Database) insertUserIntoDB(username string) (User, error) {
	var user User
	if username == "" {
		return user, fmt.Errorf("error! Missing username")
	}
	query := `INSERT INTO users (username, count) VALUES (?, ?)`
	switch db.dbType {
	case postgres:
		query = `INSERT INTO users (username, count) VALUES ($1, $2)`
	}
	err := executeQuery(db.db, query, username, "0")
	if err != nil {
		log.Panic(err)
		return user, err
	}
	user.username = username
	user.counter = 0
	return user, nil
}

func (db *Database) updateUserViewCount(user User) error {
	user.counter += 1
	query := `UPDATE users SET count=? WHERE username=?`
	switch db.dbType {
	case postgres:
		query = `UPDATE users SET count=$1 WHERE username=$2`
	}
	return executeQuery(db.db, query, fmt.Sprint(user.counter), user.username)
}

func (db *Database) searchForUser(username string) (User, error) {
	var user User
	query := `SELECT * FROM users WHERE username=?`
	switch db.dbType {
	case postgres:
		query = `SELECT * FROM users WHERE username=$1`
	}
	stmt, err := db.db.Prepare(query)

	if err != nil {
		log.Print(err)
		return user, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(username)

	if err != nil {
		log.Print(err)
		return user, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.username, &user.counter)
		if err != nil {
			log.Print(err)
			return user, err
		}
	}

	if user.username == "" {
		return user, fmt.Errorf(`error! User "%s" not found`, username)
	}

	return user, nil
}

func (db *Database) Close() {
	db.db.Close()
}
