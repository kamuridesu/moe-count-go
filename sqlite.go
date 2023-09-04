package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func StartDB(dbFileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFileName)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS "users" (
		"username"  TEXT NOT NULL UNIQUE,
		"count"     INTEGER NOT NULL,
		PRIMARY KEY("username")
	);`

	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}
	return db, nil
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

func insertUserIntoDB(db *sql.DB, username string) (User, error) {
	var user User
	if username == "" {
		return user, fmt.Errorf("error! Missing username")
	}
	err := executeQuery(db, `INSERT INTO "users" (username, count) VALUES (?, ?)`, username, "0")
	if err != nil {
		log.Print(err)
		return user, err
	}
	user.username = username
	user.counter = 0
	return user, nil
}

func updateUserViewCount(db *sql.DB, user User) error {
	user.counter += 1
	return executeQuery(db, `UPDATE users SET count=? WHERE username=?`, fmt.Sprint(user.counter), user.username)
}

func searchForUser(db *sql.DB, username string) (User, error) {
	var user User
	stmt, err := db.Prepare(`SELECT * FROM users WHERE username=?`)

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
