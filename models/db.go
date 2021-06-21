package models

import "errors"

type DB interface {
	Init()
	CheckLoginDetails(username string, password string) (uint64, error)
}

type TempDb struct {
}

var user = User{
	ID:       1,
	Username: "kangsan",
	Password: "1234",
}

func (t *TempDb) Init() {
}

func (t *TempDb) CheckLoginDetails(username string, password string) (userId uint64, err error) {
	if username != user.Username || password != user.Password {
		return 0, errors.New("Invalid login details")
	}

	//실제 db에서는 유저 id 반환
	return 1, nil
}
