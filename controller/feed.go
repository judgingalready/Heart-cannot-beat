package controller

import (
	// "fmt"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	token := c.Query("token")
	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "tokenis invalid"},
		})
		return
	}

	var videos []Video
	SearchVideoForFeed(&videos)
	for i, _ := range videos {
		var count int64
		db.Model(&Like{}).
			Where("user_id = ? AND video_id = ?", user.Id, videos[i].Id).
			Count(&count)
		fmt.Println("count %d", count)
		if count > 0 {
			// 视频已被点过赞
			videos[i].IsFavorite = true
		} else {
			// 视频未被点过赞
			videos[i].IsFavorite = false
		}
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
