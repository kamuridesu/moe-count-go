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

	db, err := StartDB(dbFileName)
	defer db.Close()

	if err != nil {
		t.Errorf("Expected for no error, got %v", err)
	}

	if _, err := os.Stat(dbFileName); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected file to be created!")
	}
}

func TestInsertUserIntoDB_EmptyUsername(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := StartDB(dbFileName)
	defer db.Close()

	username := ""

	_, err := insertUserIntoDB(db, username)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInsertUserIntoDB_OneUser(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := StartDB(dbFileName)
	defer db.Close()

	username := "test"

	user, err := insertUserIntoDB(db, username)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedUser := User{"test", 0}

	if !user.equals(&expectedUser) {
		t.Errorf("Expected %v, got %v", expectedUser, user)
	}
}

func TestUpdateUserViewCount(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := StartDB(dbFileName)
	defer db.Close()

	user := User{"test", 0}

	err := updateUserViewCount(db, user)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSearchForUser(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := StartDB(dbFileName)
	defer db.Close()

	expectedUser, _ := insertUserIntoDB(db, "test")

	user, err := searchForUser(db, "test")

	if err != nil {
		t.Errorf("Expected for no error, got %v", err)
	}

	if !user.equals(&expectedUser) {
		t.Errorf("Expected %v, got %v", expectedUser, user)
	}
}

func TestSearchForUser_EmptyUsername(t *testing.T) {
	basePath := "./test_sqlite"
	defer os.RemoveAll(basePath)
	os.Mkdir(basePath, 0755)
	dbFileName := filepath.Join(basePath, "test.db")

	db, _ := StartDB(dbFileName)
	defer db.Close()

	_, err := searchForUser(db, "")

	if err == nil {
		t.Errorf("Expected for error, got nil")
	}
}
