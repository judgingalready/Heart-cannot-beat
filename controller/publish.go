package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	// 视频保存的本地路径：用户Id_时间戳_filename.mp4
	currentTime := time.Now().Unix()
	filename = fmt.Sprintf("%d_%d_%s", user.Id, currentTime, filename)
	// 视频存入Videos数据库的url：IP:Port/static/视频本地路径
	ip_port := "114.212.85.230:8080" //暂时写死
	videoName := fmt.Sprintf("%s%s", "http://"+ip_port+"/static/", filename)
	saveFile := filepath.Join("./public/", filename)

	//保存视频到本地
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//获取视频封面
	frameNum := 1
	covelBuf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(saveFile).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(covelBuf, os.Stdout).
		Run()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//保存封面到本地
	covelImage, err := imaging.Decode(covelBuf)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	covelName := saveFile[:len(saveFile)-3] + "png"
	if err := imaging.Save(covelImage, covelName); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//在Videos数据库中插入数据
	CreateVideoErr := db.Create(&Video{AuthorID: user.Id, PlayUrl: videoName, CoverUrl: videoName[:len(videoName)-3] + "png", PublishTime: time.Now().Unix()}).Error
	if CreateVideoErr != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  CreateVideoErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  filename + " 视频已上传成功！",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user_id is invalid"},
		})
		return
	}

	token := c.Query("token")
	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	var videos []Video
	SearchVideoForPublishList(user_id, &videos)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
