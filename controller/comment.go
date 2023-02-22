package controller

import (
	"fmt"
	"net/http"
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
	actionType := c.Query("action_type")

	var user User
	if userExitErr := db.Where("name = ?", token).Take(&user).Error; userExitErr == nil {
		if actionType == "1" {
			text := c.Query("comment_text")
			currentTime := fmt.Sprintf("%d-%d", time.Now().Month(), time.Now().Day())

			comment := Comment{User: user, Content: text, CreateDate: currentTime}

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
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
