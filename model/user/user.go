package user

import (
	"TinyTolk/config"
	"TinyTolk/response/user"
	"github.com/jinzhu/gorm"
	"log"
)

//用户信息表
type User struct {
	gorm.Model
	Username string `gorm:"UNIQUE" ` //账号名
	Password string `gorm:"NOT NULL"`
	Name string //用户名
	FollowCount int64 `gorm:"DEFAULT:0"`
	FollowerCount int64 `gorm:"DEFAULT:0"`
}

type UserInfo struct {
	gorm.Model
	UserId uint
	User  User `gorm:"ForeignKey:UserId"`
	Avatar string //用户图像
	BackgroudImage string //用户个人主页顶部大图
	Signature string     //用户个人简介
	TotalFavorited int64 ` gorm:"default:0"`//获赞总数
	WorkCount int64 `gorm:"default:0"`//作品数量
	FavoriteCount int64 `gorm:"default:0"`//点赞数量
}
//用户关注表

//生成用户数据表

func CreateUsersTable()error{
	db := config.DB.AutoMigrate(&User{})
	return db.Error

}
func CreateUserInfoTable()error{
	db := config.DB.AutoMigrate(&UserInfo{})
	return db.Error

}

func InsertUser(user *User)*gorm.DB{
	return config.DB.Create(user)
}

func GetUserByUsername(user *User, username string)bool{
	config.DB.First(user,"username = ?",username)
	return user.ID != 0
}

func GetUserById(user *user.User,id uint)bool{
	result := config.DB.Model(&User{}).Find(user, "id = ?", id)
	log.Println(result)
	return user.ID != 0
}

func InsertUserInfo(userInfo *UserInfo)error{
	result := config.DB.Create(userInfo)
	return result.Error
}

func GetUserInfoByUserId( user *UserInfo,userId uint)bool{
	result :=config.DB.Model(&UserInfo{}).Where("user_id = ?",userId).Find(user)
	return result.RowsAffected != 0
}
//关注表相关操作

func AddWorkCount(userId uint) error {
	result := config.DB.Model(&UserInfo{}).Where("user_id = ?", userId).UpdateColumn("work_count",gorm.Expr("work_count + ?", 1))
	return result.Error
}
func SubWorkCount(userId uint) error {
	result := config.DB.Model(&UserInfo{}).Where("user_id = ?", userId).UpdateColumn("work_count",gorm.Expr("work_count - ?", 1))
	return result.Error
}

func AddFavoriteCount(userId uint) error {
	result := config.DB.Model(&UserInfo{}).Where("user_id = ?", userId).UpdateColumn("favorite_count",gorm.Expr("work_count + ?", 1))
	return result.Error
}
func SubFavoriteCount(userId uint) error {
	result := config.DB.Model(&UserInfo{}).Where("user_id = ?", userId).UpdateColumn("favorite_count",gorm.Expr("work_count - ?", 1))
	return result.Error
}
func AddGetFavoriteCount(userId uint) error {
	result := config.DB.Model(&UserInfo{}).Where("user_id = ?", userId).UpdateColumn("total_favorited",gorm.Expr("work_count + ?", 1))
	return result.Error
}
func SubGetFavoriteCount(userId uint) error {
	result := config.DB.Model(&UserInfo{}).Where("user_id = ?", userId).UpdateColumn("total_favorited",gorm.Expr("work_count - ?", 1))
	return result.Error
}
