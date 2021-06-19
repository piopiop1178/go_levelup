package models

//각 항목에 bind tag 추가 -> json 형태로 들어왔을 때 binding 되는 key 추가
//binding:"required"일 경우 request에 해당하는 데이터가 안넘어 왔을 때 에러 발생
type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
