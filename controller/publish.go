package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	var user User
	// token先设为用户名
	userExitErr := db.Where("name = ?", token).Take(&user).Error

	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
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
	currentTime := time.Now().Unix()
	// 视频存入Videos数据库的url：IP:Port/static/视频本地路径
	ip_port := "114.212.85.230:8080" //暂时写死
	filename = fmt.Sprintf("%d_%d_%s", user.Id, currentTime, filename)
	finalName := fmt.Sprintf("%s%s", "http://"+ip_port+"/static/", filename)
	// 视频保存的本地路径：用户Id_视频描述_时间戳
	saveFile := filepath.Join("./public/", filename)

	CreateVideoErr := db.Create(&Video{AuthorID: user.Id, PlayUrl: finalName}).Error
	if CreateVideoErr != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  CreateVideoErr.Error(),
		})
		return
	}

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.SaveUploadedFile(data, saveFile)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
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
	var videos []Video
	user_id, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		panic(err)
	}
	SearchVideoForPublishList(user_id, &videos)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
