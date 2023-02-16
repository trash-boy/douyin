package Video

import (
	"TinyTolk/config"
	user3 "TinyTolk/model/user"
	"TinyTolk/response/video"
	"github.com/jinzhu/gorm"
	"time"
)

type Video struct {
	gorm.Model
	UserId uint
	User  user3.User `gorm:"ForeignKey:UserId"`
	PlayUrl string
	CoverUrl string
	FavoriteCount int64
	CommentCount int64
	Title string
}

func CreateVideoTable()error{
	db := config.DB.AutoMigrate(&Video{})
	return db.Error
}

func InsertVideo(userId uint, playUrl , coverUrl ,title string)error{
	var video Video
	video.UserId = userId
	video.PlayUrl = playUrl
	video.CoverUrl = coverUrl
	video.Title = title

	result :=config.DB.Create(&video)
	return result.Error
}

func GetVideoListByUserId(userId uint, videoList *[]video.Video)error{
	*videoList = make([]video.Video,0)

	result := config.DB.Model(&Video{UserId: userId}).Find(videoList)
	return result.Error
}

func GetVideoListByTimeStamp(timeStamp uint, videoList *[]video.Video)error{
	//将时间戳变为时间类型
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(int64(timeStamp), 0).Format(timeLayout)
	*videoList = make([]video.Video,0)
	result := config.DB.Model(&Video{}).Where("created_at <= ?" ,timeStr).Order("created_at desc").Limit(5).Find(videoList)
	return result.Error
}




