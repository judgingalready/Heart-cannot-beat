package controller

import (
	"fmt"
	"net/http"
	"strconv"
	// "sync/atomic"
	// "time"

	"github.com/gin-gonic/gin"
)

var tempChat = map[string][]Message{}

var messageIdSequence = int64(-1)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	var user User
	// token先设为用户名
	userExitErr := db.Where("name = ?", token).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// "to_user_id"转换为int64类型
	toUserIdInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id非法"})
		return
	}
	// 查找对方用户
	var toUser User
	toUserExitErr := db.Where("id = ?", toUserIdInt).Take(&toUser).Error
	if toUserExitErr == nil {
		CreateMessageErr := db.Create(&Message{ToUserId: toUserIdInt, FromUserId: user.Id, Content: content}).Error
		if CreateMessageErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "发送消息失败！"})
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 0})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	var user User
	// token先设为用户名
	userExitErr := db.Where("name = ?", token).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// "to_user_id"转换为int64类型
	toUserIdInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id非法"})
		return
	}
	// 查找对方用户
	var toUser User
	toUserExitErr := db.Where("id = ?", toUserIdInt).Take(&toUser).Error
	if toUserExitErr == nil {
		// // 查找与双方有关的消息
		// var messages []Message
		// messagesFindErr := db.Where("((to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)) AND Id > ?",
		// 	toUserIdInt, user.Id, user.Id, toUserIdInt, messageIdSequence).Order("Id").Find(&messages).Error // 暂时按照主键降序查找

		// 查找对方发来的消息
		var messages []Message
		messagesFindErr := db.Where("to_user_id = ? AND from_user_id = ? AND Id > ?",
			user.Id, toUserIdInt, messageIdSequence).Order("Id").Find(&messages).Error // 暂时按照主键降序查找
		if messagesFindErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "查找messages数据库失败"})
			return
		}
		if len(messages) > 0 {
			messageIdSequence = messages[len(messages)-1].Id
		}
		fmt.Println(messages)
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messages})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}
