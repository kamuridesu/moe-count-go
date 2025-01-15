package main

import (
	"testing"
)

func TestUser_Equals(t *testing.T) {
	user1 := User{username: "test", counter: 1}
	user2 := User{username: "test", counter: 1}
	result := user1.Equals(&user2)
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestUser_Equals_False(t *testing.T) {
	user1 := User{username: "test", counter: 1}
	user2 := User{username: "test", counter: 2}
	result := user1.Equals(&user2)
	if result {
		t.Errorf("Expected false, got true")
	}
}
