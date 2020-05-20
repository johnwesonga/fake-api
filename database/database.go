package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type (
	user struct {
		Name  string `json:"name"`
		Sound string `json:"sound"`
	}
)

var users []user

func init() {
	wd, err := os.Getwd()
	if err != nil {
		er(err)
	}

	file, err := os.Open(fmt.Sprintf("%v/database/data.json", wd))
	if err != nil {
		er(err)
	}

	byteFile, _ := ioutil.ReadAll(file)
	if err != nil {
		er(err)
	}

	err = json.Unmarshal(byteFile, &users)
	if err != nil {
		er(err)
	}
}

// GetUser returns a user.
func GetUser(name string) (*user, error) {
	for _, v := range users {
		if name == v.Name {
			return &v, nil
		}
	}
	return nil, errors.New("No user found")
}

func GetAllUsers() []user {
	return users
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
