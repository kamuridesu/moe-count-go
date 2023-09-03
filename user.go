package main

type User struct {
	username string
	counter  int
}

func (u *User) equals(user *User) bool {
	if (*u).username == (*user).username &&
		(*u).counter == (*user).counter {
		return true
	}
	return false
}
