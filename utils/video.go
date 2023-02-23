package utils

import (
	"TinyTolk/response/video"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

func FormVideoActionResponse(statusCode int32, statusMsg string)*video.VideoActionResponse{
	var video video.VideoActionResponse
	video.StatusCode = statusCode
	video.StatusMsg = statusMsg
	return &video
}

func FormVideoListResponse(statusCode int32, statusMsg string, videoList *[]video.Video)*video.VideoListResponse{
	var video video.VideoListResponse
	video.StatusCode = statusCode
	video.StatusMsg = statusMsg
	video.VideoList = *videoList
	return &video
}

func FormVideoFeedResponse(statusCode int32, statusMsg string, videoList *[]video.Video,nextTime int64)*video.VideoFeedResponse{
	var video video.VideoFeedResponse
	video.StatusCode = statusCode
	video.StatusMsg = statusMsg
	video.VideoList = *videoList
	video.NextTime = nextTime
	return &video
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			//fmt.Printf("mkdir failed![%v]\n", err)
			log.Printf("mkdir failed! [%v],", err)
		} else {
			return true, nil
		}
	}
	return false, err

}


func WriteFile(filePath string,content multipart.File) error {

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_,err = io.Copy(file, content)
	if err != nil {
		return err
	}
	return nil
}

// GetSnapshot 生成视频缩略图并保存（作为封面）
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}

	err = imaging.Save(img, snapshotPath+".jpeg")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}

	// 成功则返回生成的缩略图名
	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".jpeg"
	return
}


