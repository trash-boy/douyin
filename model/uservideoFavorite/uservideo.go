package uservideoFavorite

import "TinyTolk/config"

type UserVideoFavorite struct {
	UserId uint `gorm:"user_id; primaryKey"`
	VideoId uint `gorm:"video_id; primaryKey"`
	Status int `gorm:"status; default:0"`

}

func CreateUserVideoFavoriteTable(){
	config.DB.AutoMigrate(UserVideoFavorite{})
}

func CreateAndUpdate(userId,videoId uint, status int)error{
	if !IsExist(userId,videoId){
		return Insert(userId,videoId , status)
	}
	return  UpdateStatus(userId,videoId,status)
}
func GetVideoIdByUserId(userId uint)([]UserVideoFavorite, error){
	var userVideo []UserVideoFavorite
	result := config.DB.Select("video_id").Where("user_id = ? and status = ?",userId, 1).Find(&userVideo)
	return userVideo,result.Error
}
func IsExist(userId,videoId uint)(exist bool){
	var uservideo UserVideoFavorite
	uservideo.UserId = userId
	uservideo.VideoId = videoId
	result := config.DB.Take(&uservideo)
	return result.RowsAffected == 1
}

func Insert(userId,videoId uint, status int)error{
	var uservideo UserVideoFavorite
	uservideo.UserId = userId
	uservideo.VideoId = videoId
	uservideo.Status = status
	result := config.DB.Create(&uservideo)
	return result.Error
}

func UpdateStatus(userId,videoId uint, status int)error{
	var uservideo UserVideoFavorite
	uservideo.UserId = userId
	uservideo.VideoId = videoId
	result := config.DB.Model(&uservideo).Update("status",status)
	return result.Error
}