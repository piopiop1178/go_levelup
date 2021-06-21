package models

type DB interface {
	Init()
	CheckLoginDetails(username string, password string) bool
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

func (t *TempDb) CheckLoginDetails(username string, password string) bool {
	if username != user.Username || password != user.Password {
		return false
	}
	return true
}
