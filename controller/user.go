package controller

import (
	"fmt"
	"net/http"

	// "sync/atomic"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

// 原始用户登录map
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token, tokenErr := GenerateToken(username)
	if tokenErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token generate error"},
		})
		return
	}

	var account Account
	accountExistErr := db.Where("name = ?", username).Take(&account).Error

	if accountExistErr == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已经存在!"},
		})
		return
	} else {
		tx := db.Begin()

		if err := tx.Create(&User{Name: username}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Error creating User"},
			})
			return
		}
		if err := tx.Create(&Account{Name: username, Password: password}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Error creating count"},
			})
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			// 如果提交失败，则回滚事务
			tx.Rollback()
			fmt.Println("Error committing transaction:", err)
			return
		}

		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token, tokenErr := GenerateToken(username)
	if tokenErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token generate error"},
		})
		return
	}

	var account Account
	accountExitErr := db.Where("name = ? AND password = ?", username, password).Take(&account).Error
	if accountExitErr == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   account.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在或者密码错误！"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token paraser error"},
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     user,
	})
}
