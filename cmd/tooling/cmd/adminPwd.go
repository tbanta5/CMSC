package cmd

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateAdminPassword(pwd string) ([]byte, error) {
	// Read the password from command line.
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
