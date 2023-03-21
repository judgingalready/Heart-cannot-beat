package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type FriendUser struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty" default:"true"`
	FriendMessage string `json:"message,omitempty"`
	MsgType       int64  `json:"msgType,omitempty"`
}

type FriendUserListResponse struct {
	Response
	FriendUserList []FriendUser `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	var user User
	verifyErr := VerifyToken(token, &user)
	if verifyErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token解析错误!"},
		})
		return
	}

	if verifyErr == nil {
		// 对方用户id
		toUserId := c.Query("to_user_id")

		var toUser User
		toUserExistErr := db.Where("id = ?", toUserId).Take(&toUser).Error
		if toUserExistErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
			return
		}
		action_type := c.Query("action_type")

		// 关注操作
		if action_type == "1" {
			// 注意不能关注自己
			if user.Id == toUser.Id {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "您不能关注自己!"})
				return
			}
			// 注意不能重复关注
			var relation Relation
			relationTakeErr := db.Where("follow = ? AND follower = ?", toUserId, user.Id).Take(&relation).Error
			if relationTakeErr == nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "您不能重复关注该用户！"})
				return
			}

			// 发布关注消息
			sb := strings.Builder{}
			userIdString := strconv.FormatInt(int64(user.Id), 10)
			toUserIdString := strconv.FormatInt(int64(toUser.Id), 10)
			sb.WriteString(userIdString)
			sb.WriteString(" ")
			sb.WriteString(toUserIdString)
			FollowPublish("FollowAdd", sb.String())

			// 删除对应的redis缓存项
			_, redisErr := rdbFollow.Del(redisContext, userIdString).Result()
			if redisErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Redis数据库无法访问"})
				return
			}
			_, redisErr = rdbFollower.Del(redisContext, toUserIdString).Result()
			if redisErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Redis数据库无法访问"})
				return
			}
		} else if action_type == "2" {
			// 发布取关消息
			sb := strings.Builder{}
			userIdString := strconv.FormatInt(int64(user.Id), 10)
			toUserIdString := strconv.FormatInt(int64(toUser.Id), 10)
			sb.WriteString(userIdString)
			sb.WriteString(" ")
			sb.WriteString(toUserIdString)
			FollowPublish("FollowDel", sb.String())

			// 删除对应的redis缓存项
			_, redisErr := rdbFollower.Del(redisContext, userIdString).Result()
			if redisErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Redis数据库无法访问"})
				return
			}
			_, redisErr = rdbFollow.Del(redisContext, toUserIdString).Result()
			if redisErr != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Redis数据库无法访问"})
				return
			}
		} else {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作类型码不存在！"})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	FollowListAction(c, true) // 显示关注列表
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	FollowListAction(c, false) // 显示粉丝列表
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	user_id := c.Query("user_id")
	// 查找用户
	var user User
	userExitErr := db.Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	// 在Relation数据查找和用户相关的记录（只查找一次）
	relations := []Relation{}
	relationsFinderr := db.Where("follow = ? OR follower = ?", user_id, user_id).Find(&relations).Error
	if relationsFinderr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	// 提取用户的关注者id和粉丝id
	follows, followers := []int64{}, []int64{}
	userIdInt, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户id非法"})
		return
	}
	for _, relation := range relations {
		if relation.Follow == userIdInt {
			followers = append(followers, relation.Follower)
		} else if relation.Follower == userIdInt {
			follows = append(follows, relation.Follow)
		}
	}
	// 寻找关注者和粉丝的共同者id
	friends := []int64{}
	for i, _ := range follows {
		for j, _ := range followers {
			if follows[i] == followers[j] {
				friends = append(friends, follows[i])
				break
			}
		}
	}
	// 根据id在User数据库查找好友
	users := []User{}
	friendsFindErr := db.Where("id IN ?", friends).Find(&users).Error
	if friendsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 将Relation数据库中用户为followers的MessageId字段全置为-1
	setRelations := []Relation{}
	setRelationsFindErr := db.Where("follower = ?", user_id).Find(&setRelations).Error
	if setRelationsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
		return
	}
	if len(setRelations) > 0 {
		for i, _ := range setRelations {
			setRelations[i].MessageId = -1
		}
		setRelationsSaveErr := db.Save(&setRelations).Error
		if setRelationsSaveErr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库更新失败"})
			return
		}
	}
	// 查找与好友双方有关的最新消息
	var message Message
	friendUsers := []FriendUser{}
	for i, _ := range users {
		messageFindErr := db.Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)",
			userIdInt, users[i].Id, users[i].Id, userIdInt).Last(&message).Error // 暂时按照主键降序查找
		if messageFindErr == nil {
			msgType := 0
			if message.FromUserId == userIdInt {
				msgType = 1
			}
			friendUsers = append(friendUsers, FriendUser{Id: users[i].Id, Name: users[i].Name, FollowCount: users[i].FollowCount,
				FollowerCount: users[i].FollowerCount, IsFollow: true, FriendMessage: message.Content, MsgType: int64(msgType)}) // 把好友的IsFollow设为true（已关注）
		} else {
			friendUsers = append(friendUsers, FriendUser{Id: users[i].Id, Name: users[i].Name, FollowCount: users[i].FollowCount,
				FollowerCount: users[i].FollowerCount, IsFollow: true}) // 把好友的IsFollow设为true（已关注）
		}
	}

	c.JSON(http.StatusOK, FriendUserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		FriendUserList: friendUsers,
	})
}

func FollowListAction(c *gin.Context, followType bool) {
	user_id := c.Query("user_id")
	// 查找用户
	var user User
	userExitErr := db.Where("id = ?", user_id).Take(&user).Error
	if userExitErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var rdb *redis.Client
	if followType {
		rdb = rdbFollow
	} else {
		rdb = rdbFollower
	}

	// 查找redis缓存
	redisValue, redisErr := rdb.SMembers(redisContext, user_id).Result()
	if redisErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Redis数据库无法访问"})
		return
	}
	follows := []int64{}
	// 缓存未命中
	if len(redisValue) == 0 {
		// 在Relation数据查找关注者的id
		relations := []Relation{}
		var relationsFinderr error
		if followType {
			relationsFinderr = db.Where("follower = ?", user_id).Find(&relations).Error
		} else {
			relationsFinderr = db.Where("follow = ?", user_id).Find(&relations).Error
		}
		if relationsFinderr != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
			return
		}
		// 如果为空直接返回
		if len(relations) == 0 {
			c.JSON(http.StatusOK, UserListResponse{
				Response: Response{
					StatusCode: 0,
				},
				UserList: []User{},
			})
			return
		}
		for _, relation := range relations {
			if followType {
				follows = append(follows, relation.Follow)
			} else {
				follows = append(follows, relation.Follower)
			}
		}
		// 更新缓存（5秒过期）
		followsString := []string{}
		for _, follow := range follows {
			followsString = append(followsString, strconv.FormatInt(int64(follow), 10))
		}
		rdb.SAdd(redisContext, user_id, followsString)
		rdb.Expire(redisContext, user_id, time.Second*5).Result()
	} else { //缓存命中
		redisFollows, err := rdb.SMembers(redisContext, user_id).Result()
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Relations数据库查询失败"})
			return
		}
		rdb.Expire(redisContext, user_id, time.Second*5).Result() // 更新缓存时间
		for _, redisValue := range redisFollows {
			redisValueInt, _ := strconv.ParseInt(redisValue, 10, 64)
			follows = append(follows, redisValueInt)
		}
	}

	// 根据id在User数据库查找
	users := []User{}
	followsFindErr := db.Where("id IN ?", follows).Find(&users).Error
	if followsFindErr != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库查询失败"})
		return
	}
	// 把users的IsFollow设为true或者false（已关注或未关注）
	for i, _ := range users {
		if followType {
			users[i].IsFollow = true
		} else {
			users[i].IsFollow = false
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}
