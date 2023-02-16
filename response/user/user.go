package user

type UserRegisterResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	UserId uint `json:"user_id"`
	Token string `json:"token"`

}


type UserLoginResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	UserId uint `json:"user_id"`
	Token string `json:"token"`
}


type User struct {
	ID uint `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"`
	FollowCount int64 `json:"follow_count" gorm:"follow_count"`
	FollowerCount int64 `json:"follower_count" gorm:"follower_count"`
	IsFollow bool `json:"is_follow" gorm:"is_follow"`
}
type UserResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	User User `json:"user"`
}

