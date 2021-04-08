package game

import (
	"fmt"
	"strings"

	"github.com/labiraus/gomud-common/db"
)

func (u *User) login() (bool, error) {
	<-u.Send("Username:")
	username, err := u.read()
	if err != nil {
		return false, err
	}

	<-u.Send("Password:")
	password, err := u.read()
	if err != nil {
		return false, err
	}
	user := db.Login(username[0], strings.Join(password, " "))
	u.id = user.ID
	u.name = user.CharacterName
	if u.id == 0 {
		<-u.Send("Login failed")
		return false, nil
	}

	return true, nil
}

func (u *User) quit() (bool, error) {
	<-u.Send("Are you sure? y/N")
	command, err := u.readBlank()
	if err != nil {
		return false, err
	}
	return strings.ToLower(command[0]) == "y", nil
}

func (u *User) createLogin() error {
	username, err := u.createUserName()
	if err != nil {
		return err
	}

	password, err := u.createPassword()
	if err != nil {
		return err
	}
	characterName, err := u.createCharacter()
	if err != nil {
		return err
	}

	u2 := db.CreateUser(username, password, characterName)
	fmt.Printf("%+v", u2)
	u.id = u2.ID
	u.name = characterName
	return nil
}

func (u *User) createUserName() (string, error) {
	for {
		name := ""
		<-u.Send("Username:")

		for name == "" {
			nameArray, err := u.readBlank()
			if err != nil {
				return "", err
			}
			name = nameArray[0]
		}
		nameExists := db.UserNameExists(name)
		if !nameExists {
			return name, nil
		}
		<-u.Send("Username not available")
	}
}

func (u *User) createCharacter() (string, error) {
	for {
		name := ""
		<-u.Send("Character name:")

		for name == "" {
			nameArray, err := u.readBlank()
			if err != nil {
				return "", err
			}
			name = nameArray[0]
		}
		nameExists := db.CharacterNameExists(name)
		if !nameExists {
			return name, nil
		}
		<-u.Send("Character name not available")
	}
}

func (u *User) createPassword() (string, error) {
	for {
		password := ""
		<-u.Send("Password:")

		for password == "" {
			nameArray, err := u.readBlank()
			if err != nil {
				return "", err
			}
			password = nameArray[0]
		}

		if len(password) > 6 {
			return password, nil
		}
		<-u.Send("Password does not meet complexity requirements")
	}
}
