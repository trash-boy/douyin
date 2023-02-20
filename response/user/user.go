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
	Avatar string `json:"avatar" gorm:"avatar"`//用户图像
	BackgroudImage string `json:"backgroud_image" gorm:"backgroud_image"`//用户个人主页顶部大图
	Signature string     `json:"signature" gorm:"signature"`//用户个人简介
	TotalFavorited int64 `json:"total_favorited" gorm:"total_favorited,default:0"`//获赞总数
	WorkCount int64 `json:"work_count" gorm:"work_count,default:0"`//作品数量
	FavoriteCount int64 `json:"favorite_count" gorm:"favorite_count,default:0"`//点赞数量
}
type UserResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	User User `json:"user"`
}

