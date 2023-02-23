package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
	}

	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	if actionType == "1" {
		text := c.Query("comment_text")
		currentTime := fmt.Sprintf("%d-%d", time.Now().Month(), time.Now().Day())

		comment := Comment{User: user, Content: text, CreateDate: currentTime, VideoId: videoId}

		createCommentErr := db.Create(&comment).Error
		if createCommentErr != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  createCommentErr.Error(),
			})
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: Comment{
				Id:         comment.Id,
				User:       user,
				Content:    text,
				CreateDate: currentTime,
			}})
		return
	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		deleteCommentErr := db.Where("id = ?", commentId).Delete(&Comment{}).Error
		if deleteCommentErr != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  deleteCommentErr.Error(),
			})
		}
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	token := c.Query("token")
	videoId := c.Query("video_id")

	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "tokenis invalid"},
		})
		return
	}

	var comments []Comment
	err := db.Where("comments.video_id = ?", videoId).Order("comments.id desc").
		Find(&comments).Error

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "cannot get comments"},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
		},
		CommentList: comments,
	})

}
