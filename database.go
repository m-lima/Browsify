package main

import (
	"errors"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/markbates/goth"
)

type User struct {
	gorm.Model
	Email               string `gorm:"not null;unique"`
	Avatar              string
	Admin               bool
	CanShowHidden       bool
	ShouldShowHidden    bool
	CanShowProtected    bool
	ShouldShowProtected bool
}

var (
	dbLogStd = log.New(os.Stdout, "[db] ", log.Ldate|log.Ltime)
	dbLogErr = log.New(os.Stderr, "ERROR [db] ", log.Ldate|log.Ltime)

	db *gorm.DB
)

func buildConnectionString() string {
	connectionString := ""
	if Configuration.Db.Host != "" {
		connectionString += " host=" + Configuration.Db.Host
	}
	if Configuration.Db.Name != "" {
		connectionString += " dbname=" + Configuration.Db.Name
	}
	if Configuration.Db.User != "" {
		connectionString += " user=" + Configuration.Db.User
	}
	if Configuration.Db.Password != "" {
		connectionString += " password=" + Configuration.Db.Password
	}
	if Configuration.Db.SslMode != "" {
		connectionString += " sslmode=" + Configuration.Db.SslMode
	}

	return connectionString
}

func InitDB() error {
	var err error
	db, err = gorm.Open("postgres", buildConnectionString())
	if err != nil {
		dbLogErr.Println("failed to connect to database")
		return err
	}

	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		dbLogErr.Println("failed to automigrate user")
		return err
	}

	dbLogStd.Println("successful DB initialization")

	return nil
}

func FinalizeDB() {
	db.Close()
	dbLogStd.Println("successful DB finalization")
}

func ValidateUser(sessionUser *goth.User) bool {
	if session.Email == "" {
		return false
	}

	dbLogStd.Println("user logging in:", sessionUser.Email)

	var user User

	if !db.Where("email = ?", sessionUser.Email).First(&user).RecordNotFound() {
		err := db.Model(&user).Update("avatar", sessionUser.AvatarURL).Error
		if err == nil {
			dbLogStd.Println("user avatar updated:", user.Email)
		} else {
			dbLogErr.Println("user avatar not updated:", user.Email, err)
		}
	}

	return true
}

func CreateUser(newUser *User) error {
	if db.NewRecord(newUser) {
		err := db.Create(newUser).Error
		if err != nil {
			dbLogErr.Println("could not create user:", newUser.Email)
			return err
		}

		return nil
	}

	return errors.New("user already exits")
}

func ReadUser(sessionUser *goth.User) (User, error) {
	var user User

	if !db.Where("email = ?", sessionUser.Email).First(&user).RecordNotFound() {
		return user, nil
	}

	dbLogErr.Println("user is not authorized:", sessionUser.Email)
	return User{}, errors.New("user is not authorized:" + sessionUser.Email)
}

func UpdateUser(updatedUser *User) error {
	err := db.Save(updatedUser).Error
	if err == nil {
		dbLogStd.Println("user updated:", updatedUser.Email)
		return nil
	} else {
		dbLogErr.Println("could not update user:", updatedUser.Email, err)
		return err
	}
}
