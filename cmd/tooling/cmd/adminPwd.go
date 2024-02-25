package cmd

import (
	"flag"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdminPassword() ([]byte, error) {
	// Read the password from command line.
	pwd := flag.String("password", "password123$", "Use this to pass in a unique password")
	flag.Parse()
	hash, err := bcrypt.GenerateFromPassword([]byte(*pwd), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
