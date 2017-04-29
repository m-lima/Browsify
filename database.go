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
	LogStd = log.New(os.Stdout, "db: ", 0)
	LogErr = log.New(os.Stderr, "db: ", 0)

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
		LogErr.Println("failed to connect to database")
		return err
	}

	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		LogErr.Println("failed to automigrate user")
		return err
	}

	LogStd.Println("successful DB initialization")

	return nil
}

func FinalizeDB() {
	db.Close()
	LogStd.Println("successful DB finalization")
}

func ValidateUser(sessionUser *goth.User) bool {
	LogStd.Println("user logging in:", sessionUser.Email)

	var user User

	if !db.Where("email = ?", sessionUser.Email).First(&user).RecordNotFound() {
		err := db.Model(&user).Update("avatar", sessionUser.AvatarURL).Error
		if err == nil {
			LogStd.Println("user avatar updated:", user.Email)
		} else {
			LogErr.Println("user avatar not updated:", user.Email, err)
		}
	}

	return true
}

func CreateUser(newUser *User) error {
	if db.NewRecord(newUser) {
		err := db.Create(newUser).Error
		if err != nil {
			LogErr.Println("could not create user:", newUser.Email)
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

	LogErr.Println("user is not authorized:", sessionUser.Email)
	return User{}, errors.New("user is not authorized:" + sessionUser.Email)
}

func UpdateUser(updatedUser *User) error {
	err := db.Save(updatedUser).Error
	if err == nil {
		LogStd.Println("user updated:", updatedUser.Email)
		return nil
	} else {
		LogErr.Println("could not update user:", updatedUser.Email, err)
		return err
	}
}
