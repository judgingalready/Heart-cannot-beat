package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	var num int32
	var videos []Video
	var user User

	userExitErr := db.Where("name = ?", token).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token is invalid"},
		})
		return
	}
	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "video_id is invalid"},
		})
		return
	}

	if actionType == "1" {
		num = 1
		likes := Like{uint(user.Id), uint(video_id)}
		db.Create(likes)
	} else if actionType == "2" {
		num = -1
		db.Where("user_id = ?", uint(user.Id)).Delete(&Like{})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "actionType is invalid"},
		})
	}
	VideoForAction(video_id, &videos, num)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "点赞已上传成功！",
	})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
