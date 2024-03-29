package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64 `json:"id,omitempty" gorm:"primary_key"`
	Author        User  `json:"author" gorm:"foreignKey:Id;references:AuthorID;"`
	AuthorID      int64
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	PublishTime   int64
}

type Comment struct {
	Id         int64 `json:"id,omitempty" gorm:"primary_key;"`
	User       User  `json:"user" gorm:"foreignKey:Id;references:UserID;"`
	UserID     int64
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoId    int64
}

type User struct {
	Id             int64  `json:"id,omitempty" gorm:"primary_key"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	IsFollow       bool   `json:"is_follow,omitempty"`
	TotalFavorited int64  `json:"total_favorited,omitempty" gorm:"default:0"`
	WorkCount      int64  `json:"work_count,omitempty" gorm:"default:0"`
	FavoriteCount  int64  `json:"favorite_count,omitempty" gorm:"default:0"`
}

type Account struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type Like struct {
	UserID  uint `gorm:"userid"`
	VideoID uint `gorm:"videoid"`

	// User    User  `gorm:"foreignkey:UserID"`
	// Video   Video `gorm:"foreignkey:VideoID"`
}

type CommentForVideo struct {
	VideoID int64   `gorm:"videoid"`
	comment Comment `json:"comment"`
}

type Relation struct {
	Id        int64 `json:"id,omitempty"`
	Follow    int64 `json:"follow,omitempty"`
	Follower  int64 `json:"follower,omitempty"`
	MessageId int64 `json:"message_id,omitempty"` // 为了确认聊天是否是刚开始已经对方发的消息是否最新
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
