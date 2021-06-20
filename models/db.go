package models

type DB interface {
	Init()
	CheckLoginDetails(username string, password string) bool
	SaveTokenToDb(userid int64, ti TokenInfo) error
}

type TempDb struct {
}

var user User

func (t *TempDb) Init() {
	//뭐 넣을까 이거 되나??
	user = User{
		ID:       1,
		Username: "kangsan",
		Password: "1234",
	}
}

func (t *TempDb) CheckLoginDetails(username string, password string) bool {
	if username != user.Username || password != user.Password {
		return false
	}
	return true
}

func (t *TempDb) SaveTokenToDb() error {
	//tmp list에 저장?
	return nil
}

type Redis struct {
}
