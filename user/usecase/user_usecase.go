package usecase

import (
	"encoding/json"
	"legato_server/api"
	"legato_server/authenticate"
	legatoDb "legato_server/db"
	"legato_server/domain"
	"legato_server/helper/converter"
	"log"
	"strings"
	"time"
)

var defaultUser = legatoDb.User{
	Username: "legato",
	Email:    "legato@gmail.com",
	Password: "1234qwer",
}

type userUseCase struct {
	db             *legatoDb.LegatoDB
	contextTimeout time.Duration
}

func NewUserUseCase(db *legatoDb.LegatoDB, timeout time.Duration) domain.UserUseCase {
	return &userUseCase{
		db:             db,
		contextTimeout: timeout,
	}
}

// Register new user and add it in our database
func (u *userUseCase) RegisterNewUser(nu api.NewUser) error {
	err := u.db.AddUser(converter.NewUserToUserDb(nu))
	if err != nil {
		return err
	}

	return nil
}

// Login the user
// return the access token
func (u *userUseCase) Login(cred api.UserCredential) (t string, e error) {
	// Check username validation
	expectedUser, err := u.db.GetUserByUsername(cred.Username)
	if err != nil {
		return "", err
	}

	// Check credentials
	token, err := authenticate.Login(cred, expectedUser)
	if err != nil {
		return "", err
	}

	t = token.TokenString

	return t, nil
}

// Returns user that has the email address
func (u *userUseCase) GetUserByEmail(s string) (user api.UserInfo, e error) {
	ue, err := u.db.GetUserByEmail(s)
	user = converter.UserDbToUser(ue)
	if err != nil {
		return api.UserInfo{}, err
	}

	return user, nil
}

// Returns user that has the username
func (u *userUseCase) GetUserByUsername(s string) (user api.UserInfo, e error) {
	ue, err := u.db.GetUserByUsername(s)
	user = converter.UserDbToUser(ue)
	if err != nil {
		return api.UserInfo{}, err
	}

	return user, nil
}

// Returns a list of all of our users in database
func (u *userUseCase) GetAllUsers() (users []*api.UserInfo, e error) {
	us, err := u.db.FetchAllUsers()
	if err != nil {
		return users, err
	}

	for _, u := range us {
		user := converter.UserDbToUser(u)
		users = append(users, &user)
	}

	return users, nil
}

func (u *userUseCase) RefreshUserToken(at string) (api.RefreshToken, error) {
	t, err := authenticate.Refresh(at)
	if err != nil {
		return api.RefreshToken{}, err
	}

	return api.RefreshToken{RefreshToken: t.TokenString}, nil
}

// This is for testing purposes
// It puts default user in the database.
func (u *userUseCase) CreateDefaultUser() error {
	err := u.db.AddUser(defaultUser)
	if err != nil {
		log.Printf("Default user is not created: %v\n", err)
		return err
	}
	log.Printf("Default user is created: %v\n", defaultUser)

	return nil
}

func (u *userUseCase) AddConnectionToDB(name string, ut api.Connection) (api.Connection, error) {
	user, _ := u.db.GetUserByUsername(name)
	con := legatoDb.Connection{}
	con.Name = ut.Name
	var err error
	var jsonData map[string]interface{}
	con.Data, jsonData, err = converter.ExtractData(ut.Data, ut.Type, &ut)
	if err != nil || strings.EqualFold(con.Data, "null") {
		return api.Connection{}, err
	}
	con.UserID = uint(ut.ID)
	con.Type = ut.Type
	c, err := u.db.AddConnection(&user, con)
	if err != nil {
		return api.Connection{}, err
	}
	ut.Data = jsonData
	ut.ID = uint(c.ID)
	return ut, nil
}

func (u *userUseCase) GetConnectionByID(username string, id uint) (legatoDb.Connection, error) {
	user, _ := u.db.GetUserByUsername(username)
	connection, err := u.db.GetUserConnectionById(&user, id)
	return connection, err
}

func (u *userUseCase) GetConnections(username string) ([]legatoDb.Connection, error) {
	user, _ := u.db.GetUserByUsername(username)
	connections, err := u.db.GetUserConnections(&user)
	if err != nil {
		return nil, err
	}
	return connections, nil
}

func (u *userUseCase) UpdateUserConnectionNameById(username string, ut api.Connection) error {
	user, _ := u.db.GetUserByUsername(username)

	err := u.db.UpdateUserConnectionNameByID(&user, ut.Name, uint(ut.ID))

	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) CheckConnectionByID(username string, id uint) error {
	user, _ := u.db.GetUserByUsername(username)
	err := u.db.CheckConnectionByID(&user, id)
	if err != nil {
		return err
	}
	return err
}

func (u *userUseCase) DeleteUserConnectionById(username string, id uint) error {
	user, _ := u.db.GetUserByUsername(username)
	err := u.db.DeleteConnectionByID(&user, id)

	if err != nil {
		return err
	}

	return nil
}
func (u *userUseCase) UpdateDataConnectionByID(username string, ut api.Connection) error {
	user, _ := u.db.GetUserByUsername(username)
	jsonString, err := json.Marshal(ut.Data)
	err = u.db.UpdateDataFieldByID(&user, string(jsonString), uint(ut.ID))

	if err != nil {
		return err
	}

	return nil
}
