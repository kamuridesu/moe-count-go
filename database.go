package main

import (
	"database/sql"
	"errors"
	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseType int

const (
	sqlite DatabaseType = iota
	postgres
)

type Database interface {
	connect() error
	close() error
	createTable() error
	dropTable() error
	searchForUser(username string) (*User, error)
	insertUserIntoDB(username string) (*User, error)
	updateUserViewCount(user *User) error
	deleteUser(username string) error
}

type Query struct {
	dbType  DatabaseType
	content string
}

type Queries struct {
	createTable *Query
	dropTable   *Query
	searchUser  *Query
	insertUser  *Query
	updateUser  *Query
	deleteUser  *Query
}

type GenericDB struct {
	db      *sql.DB
	queries *Queries
}

func (q *Query) getQuery() string {
	switch q.dbType {
	case sqlite:
		return q.content
	case postgres:
		newQuery := ""
		counter := 1
		for i := 0; i < len(q.content); i++ {
			if q.content[i] == '?' {
				newQuery += "$" + strconv.Itoa(counter)
				counter++
				continue
			}
			newQuery += string(q.content[i])
		}
		return newQuery
	}
	return q.content
}

func newQuery(dbType DatabaseType, content string) *Query {
	return &Query{dbType: dbType, content: content}
}

func newQueries(dbType DatabaseType) *Queries {
	return &Queries{
		createTable: newQuery(dbType, "CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, count INTEGER)"),
		dropTable:   newQuery(dbType, "DROP TABLE IF EXISTS users"),
		searchUser:  newQuery(dbType, "SELECT username, count FROM users WHERE username = ?"),
		insertUser:  newQuery(dbType, "INSERT INTO users (username, count) VALUES (?, ?)"),
		updateUser:  newQuery(dbType, "UPDATE users SET count = ? WHERE username = ?"),
		deleteUser:  newQuery(dbType, "DELETE FROM users WHERE username = ?"),
	}
}

func (g *GenericDB) connect() error {
	return g.db.Ping()
}

func (g *GenericDB) close() error {
	return g.db.Close()
}

func (g *GenericDB) createTable() error {
	_, err := g.db.Exec(g.queries.createTable.getQuery())
	if err != nil {
		return err
	}
	return nil
}

func (g *GenericDB) dropTable() error {
	_, err := g.db.Exec(g.queries.dropTable.getQuery())
	if err != nil {
		return err
	}
	return nil
}

func (g *GenericDB) searchForUser(username string) (*User, error) {
	var user User
	row := g.db.QueryRow(g.queries.searchUser.getQuery(), username)
	err := row.Scan(&user.username, &user.counter)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (g *GenericDB) insertUserIntoDB(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("error! Missing username")
	}
	_, err := g.db.Exec(g.queries.insertUser.getQuery(), username, 0)
	if err != nil {
		return nil, err
	}
	return &User{username: username, counter: 0}, nil
}

func (g *GenericDB) updateUserViewCount(user *User) error {
	user.counter++
	_, err := g.db.Exec(g.queries.updateUser.getQuery(), user.counter, user.username)
	if err != nil {
		return err
	}
	return nil
}

func (g *GenericDB) deleteUser(username string) error {
	_, err := g.db.Exec(g.queries.deleteUser.getQuery(), username)
	if err != nil {
		return err
	}
	return nil
}

func OpenPostgresDB(parameters string) (*sql.DB, error) {
	db, err := sql.Open("postgres", parameters)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	return db, err
}

func OpenSqliteDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDatabase(dbType DatabaseType, dbParams string) (Database, error) {
	var db *sql.DB
	var err error
	switch dbType {
	case sqlite:
		db, err = OpenSqliteDB(dbParams)
	case postgres:
		db, err = OpenPostgresDB(dbParams)
	default:
		return nil, errors.New("unknown database type")
	}
	if err != nil {
		return nil, err
	}
	queries := newQueries(dbType)
	genericDB := &GenericDB{db: db, queries: queries}
	err = genericDB.connect()
	if err != nil {
		return nil, err
	}
	err = genericDB.createTable()
	if err != nil {
		return nil, err
	}
	return genericDB, err
}
