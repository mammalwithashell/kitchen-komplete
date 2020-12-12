package main

import (
	"fmt"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Hash Passwords
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check passwords
func checkHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Check Hash Err:", err)
		return false
	}
	return true
}

func auth(session *sessions.Session, err error) bool {
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		return true
	}
	return false
}
