package postgres

import (
	"errors"
	"fmt"
	"legato_server/internal/legato/database/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string
	Email       string
	Password    string
	LastLogin   time.Time
	Scenarios   []Scenario
	Webhooks    []Webhook
	Services    []Service
	Connections []Connection
}

func (u *User) model() models.User {
	return models.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		LastLogin: u.LastLogin,
	}
}

func (u *User) string() string {
	return fmt.Sprintf("(@User: %+v)", *u)
}

func (ldb *LegatoDB) AddUser(u models.User) error {
	// Encode the user password
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err != nil {
		return err
	} else {
		u.Password = string(pw)
	}

	// Check unique username
	user := &User{}
	ldb.db.Where(&User{Username: u.Username}).Find(&user)
	if user.Username == u.Username {
		return errors.New("this username is already taken")
	}

	// Check unique user email
	user = &User{}
	ldb.db.Where(&User{Email: u.Email}).Find(&user)
	if user.Email == u.Email {
		return errors.New("this email is already taken")
	}

	// Create the database model
	newUser := User{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		LastLogin: time.Now(),
	}

	ldb.db.Create(&newUser)

	return nil
}

func (ldb *LegatoDB) GetUserByUsername(username string) (models.User, error) {
	user := User{}
	ldb.db.Where(&User{Username: strings.ToLower(username)}).First(&user)
	if user.Username != username {
		return models.User{}, errors.New("username does not exist")
	}

	return user.model(), nil
}

func (ldb *LegatoDB) GetUserById(userId uint) (models.User, error) {
	user := User{}
	ldb.db.Where("id = ?", userId).First(&user)
	if user.ID != userId {
		return models.User{}, errors.New("user id does not exist")
	}

	return user.model(), nil
}

func (ldb *LegatoDB) GetUserByEmail(email string) (models.User, error) {
	user := User{}
	ldb.db.Where(&User{Email: strings.ToLower(email)}).First(&user)
	if user.Email != email {
		return models.User{}, errors.New("email does not exist")
	}

	return user.model(), nil
}

func (ldb *LegatoDB) FetchAllUsers() ([]models.User, error) {
	var users []User
	ldb.db.Find(&users)

	if len(users) <= 0 {
		return []models.User{}, errors.New("there is no user")
	}

	var userModels []models.User
	for _, u := range users {
		userModels = append(userModels, u.model())
	}

	return userModels, nil
}
