package controller

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	// token先设为用户名
	token := c.Query("token")

	var user User
	userExistErr := db.Where("name = ?", token).Take(&user).Error
	if userExistErr == nil {
		// 对方用户id
		to_user_id := c.Query("to_user_id")

		var to_user User
		to_userExistErr := db.Where("id = ?", to_user_id).Take(&to_user).Error
		if to_userExistErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
			return
		}
		action_type := c.Query("action_type")

		// 关注操作
		if action_type == "1" {
			// 注意不能关注自己
			if user.Id == to_user.Id {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "您不能关注自己"})
				return
			}

			CreateFollowErr := db.Create(&Relation{Follow: to_user.Id, Follower: user.Id}).Error
			if CreateFollowErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库插入失败"})
				return
			}

			user.FollowCount += 1
			userSaveErr := db.Save(&user).Error
			if userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}

			to_user.FollowerCount += 1
			to_userSaveErr := db.Save(&to_user).Error
			if to_userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}

		} else if action_type == "2" {
			// 取关操作
			DeleteFollowErr := db.Where("follow = ? AND follower = ?", to_user.Id, user.Id).Delete(&Relation{}).Error
			if DeleteFollowErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库删除失败"})
				return
			}

			if user.FollowCount > 0 {
				user.FollowCount -= 1
			}
			userSaveErr := db.Save(&user).Error
			if userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}

			if to_user.FollowerCount > 0 {
				to_user.FollowerCount -= 1
			}
			to_userSaveErr := db.Save(&to_user).Error
			if to_userSaveErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
				return
			}
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 查找用户
	var user User
	userExitErr := db.Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// 在Relation数据查找关注的人的id
	relations := []Relation{}
	relationsFinderr := db.Where("follower = ?", user_id).Find(&relations).Error
	if relationsFinderr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	follows := []int64{}
	for _, relation := range relations {
		follows = append(follows, relation.Follow)
	}
	// 根据id在User数据库查找
	users := []User{}
	followsFindErr := db.Where("id IN ?", follows).Find(&users).Error
	if followsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 把users的IsFollow设为true（已关注）
	for i, _ := range users {
		users[i].IsFollow = true
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 查找用户
	var user User
	userExitErr := db.Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// 在Relation数据查找粉丝的id
	relations := []Relation{}
	relationsFinderr := db.Where("follow = ?", user_id).Find(&relations).Error
	if relationsFinderr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	followers := []int64{}
	for _, relation := range relations {
		followers = append(followers, relation.Follower)
	}
	// 根据id在User数据库查找
	users := []User{}
	followersFindErr := db.Where("id IN ?", followers).Find(&users).Error
	if followersFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 把users的IsFollow设为false（未关注）
	for i, _ := range users {
		users[i].IsFollow = false
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}
