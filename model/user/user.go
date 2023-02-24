package user

import (
	"TinyTolk/config"
	"TinyTolk/utils/encryption"
	"github.com/jinzhu/gorm"
	"time"
)

// 用户信息表
type User struct {
	gorm.Model
	Username       string `gorm:"UNIQUE" json:"-" ` //账号名
	Password       string `gorm:"NOT NULL" json:"-"`
	Name           string `json:"name"` //用户名
	FollowCount    int64  `gorm:"DEFAULT:0" json:"follow_count"`
	FollowerCount  int64  `gorm:"DEFAULT:0" json:"follower_count"`
	Avatar         string //用户图像
	BackgroudImage string //用户个人主页顶部大图
	Signature      string //用户个人简介
	TotalFavorited int64  ` gorm:"default:0"` //获赞总数
	WorkCount      int64  `gorm:"default:0"`  //作品数量
	FavoriteCount  int64  `gorm:"default:0"`  //点赞数量

}

//type UserInfo struct {
//	gorm.Model
//	UserId uint
//	User  User `gorm:"ForeignKey:UserId"`
//	Avatar string //用户图像
//	BackgroudImage string //用户个人主页顶部大图
//	Signature string     //用户个人简介
//	TotalFavorited int64 ` gorm:"default:0"`//获赞总数
//	WorkCount int64 `gorm:"default:0"`//作品数量
//	FavoriteCount int64 `gorm:"default:0"`//点赞数量
//}
//用户关注表

//生成用户数据表

func CreateUsersTable() error {
	db := config.DB.AutoMigrate(&User{})
	return db.Error

}

//func CreateUserInfoTable()error{
//	db := config.DB.AutoMigrate(&UserInfo{})
//	return db.Error
//
//}

// 返回插入用户的id
func InsertUser(userName, password string) (uint, error) {
	var user User
	user.Username = userName
	user.Password = encryption.Encoding(password)
	user.Name = config.UserPrefix + userName
	user.Avatar = config.AVATAR_URL
	user.BackgroudImage = config.BACKGROUD_IMAGE
	user.CreatedAt = time.Now()
	result := config.DB.Create(&user)
	return user.ID, result.Error
}

func GetUserByUsername(user *User, username string) bool {
	config.DB.First(user, "username = ?", username)
	return user.ID != 0
}

func GetUserById(user interface{}, id uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", id).Scan(user)
	return result.Error
}

//func InsertUserInfo(userInfo *UserInfo)error{
//	result := config.DB.Create(userInfo)
//	return result.Error
//}
//
//func GetUserInfoByUserId( user *UserInfo,userId uint)bool{
//	result :=config.DB.Model(&UserInfo{}).Where("user_id = ?",userId).Find(user)
//	return result.RowsAffected != 0
//}
//关注表相关操作

func AddWorkCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("work_count", gorm.Expr("work_count + ?", 1))
	return result.Error
}
func SubWorkCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("work_count", gorm.Expr("work_count - ?", 1))
	return result.Error
}

func AddFavoriteCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	return result.Error
}
func SubFavoriteCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("favorite_count", gorm.Expr("favorite_count- ?", 1))
	return result.Error
}
func AddGetFavoriteCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("total_favorited", gorm.Expr("total_favorited+ ?", 1))
	return result.Error
}
func SubGetFavoriteCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("total_favorited", gorm.Expr("total_favorited - ?", 1))
	return result.Error
}

func AddFollowCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
	return result.Error
}
func SubFollowCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
	return result.Error
}

func AddFollowerCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
	return result.Error
}
func SubFollowerCount(userId uint) error {
	result := config.DB.Model(&User{}).Where("id = ?", userId).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
	return result.Error
}
