package main

import (
	"errors"

	"github.com/markbates/goth"
)

type User struct {
	Email               string
	Avatar              string
	Admin               bool
	CanShowHidden       bool
	ShouldShowHidden    bool
	CanShowProtected    bool
	ShouldShowProtected bool
}

var (
	authorizedUsers = [...]User{
		User{
			Email: "marcelo@telenordigital.com",
			// Avatar: "https://lh5.googleusercontent.com/-i2nXCcG77N0/AAAAAAAAAAI/AAAAAAAAAC0/d4xJpxg2mDM/photo.jpg",
			Admin: true,
		},
		// User{
		// 	Email: "marcelowind@gmail.com",
		// },
		User{
			Email: "kris@telenordigital.com",
		},
	}
)

func ValidateUser(sessionUser *goth.User) bool {
	for i, user := range authorizedUsers {
		if sessionUser.Email == user.Email {
			user.Avatar = sessionUser.AvatarURL
			authorizedUsers[i] = user
			return true
		}
	}

	return true
}

func GetUser(sessionUser *goth.User) (User, error) {
	for _, user := range authorizedUsers {
		if sessionUser.Email == user.Email {
			return user, nil
		}
	}

	return User{}, errors.New("user is not authorized")
}

func UpdateUser(updatedUser *User) error {
	for i, user := range authorizedUsers {
		if updatedUser.Email == user.Email {
			authorizedUsers[i] = *updatedUser
			return nil
		}
	}

	return errors.New("user is not authorized")
}
