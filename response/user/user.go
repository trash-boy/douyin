package user

import "TinyTolk/model/user"

type UserRegisterResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	UserId int64 `json:"user_id"`
	Token string `json:"token"`

}


type UserLoginResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	UserId int64 `json:"user_id"`
	Token string `json:"token"`
}


type User struct {

	user.User
	IsFollow bool `json:"is_follow" gorm:"is_follow"`

}

type UserResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	User User `json:"user"`
}

