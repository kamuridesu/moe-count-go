package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStartDB_CheckIfDBIsCreated(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, err := NewDatabase(sqlite, dbFileName)
	if err != nil {
		t.Errorf("Expected for no error, got %v", err)
	}

	defer db.close()

	if _, err := os.Stat(dbFileName); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected file to be created!")
	}
}

func TestInsertUserIntoDB_EmptyUsername(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := NewDatabase(sqlite, dbFileName)
	defer db.close()

	username := ""

	_, err := db.insertUserIntoDB(username)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInsertUserIntoDB_OneUser(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := NewDatabase(sqlite, dbFileName)
	defer db.close()

	username := "test"

	user, err := db.insertUserIntoDB(username)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedUser := User{username: "test", counter: 0}

	if !user.Equals(&expectedUser) {
		t.Errorf("Expected %v, got %v", expectedUser, user)
	}
}

func TestUpdateUserViewCount(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := NewDatabase(sqlite, dbFileName)
	defer db.close()

	user := &User{username: "test", counter: 0}

	err := db.updateUserViewCount(user)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSearchForUser(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := NewDatabase(sqlite, dbFileName)
	defer db.close()

	expectedUser, _ := db.insertUserIntoDB("test")

	user, err := db.searchForUser("test")

	if err != nil {
		t.Errorf("Expected for no error, got %v", err)
	}

	if !user.Equals(expectedUser) {
		t.Errorf("Expected %v, got %v", expectedUser, user)
	}
}

func TestSearchForUser_EmptyUsername(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := NewDatabase(sqlite, dbFileName)
	defer db.close()

	_, err := db.searchForUser("")

	if err == nil {
		t.Errorf("Expected for error, got nil")
	}
}
