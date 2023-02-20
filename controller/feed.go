package controller

import (
	// "fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const videoCount = 1

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var videos []Video
	SearchDataForVideo(&videos)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}

func SearchDataForVideo(videos *[]Video) {
	err := db.Preload("Author").Find(videos).Error
	if err != nil {
		panic(err)
	}
	// fmt.Println(videos)
}
