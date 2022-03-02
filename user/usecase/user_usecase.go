package usecase

import (
	"encoding/json"
	"legato_server/api"
	"legato_server/domain"
	"legato_server/helper/converter"
	legatoDb2 "legato_server/internal/legato/database/postgres"
	"strings"
	"time"
)

type userUseCase struct {
	db             *legatoDb2.LegatoDB
	contextTimeout time.Duration
}

func NewUserUseCase(db *legatoDb2.LegatoDB, timeout time.Duration) domain.UserUseCase {
	return &userUseCase{
		db:             db,
		contextTimeout: timeout,
	}
}

func (u *userUseCase) AddConnectionToDB(name string, ut api.Connection) (api.Connection, error) {
	user, _ := u.db.GetUserByUsername(name)
	con := legatoDb2.Connection{}
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

func (u *userUseCase) GetConnectionByID(username string, id uint) (legatoDb2.Connection, error) {
	user, _ := u.db.GetUserByUsername(username)
	connection, err := u.db.GetUserConnectionById(&user, id)
	return connection, err
}

func (u *userUseCase) GetConnections(username string) ([]legatoDb2.Connection, error) {
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
