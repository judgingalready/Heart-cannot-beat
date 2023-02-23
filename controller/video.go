package controller

import (
	"gorm.io/gorm"
)

const videoCount = 30

func SearchVideoForFeed(videos *[]Video, latestTime int64) {
	err := db.Preload("Author").Order("videos.publish_time desc").Where("videos.publish_time < ?", latestTime).
		Limit(videoCount).Find(videos).Error
	if err != nil {
		panic(err)
	}
}

func SearchVideoForPublishList(user_id int64, videos *[]Video) {
	err := db.Where("author_id = ?", user_id).Find(videos).Error
	if err != nil {
		panic(err)
	}
}

func VideoForAction(video_id int64, video *Video, num int32) error {
	err := db.Where("id = ?", video_id).Find(video).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", num)).Error
	favorateErr := db.Model(&User{}).
		Where("id = ?", video.AuthorID).
		UpdateColumn("total_favorated", gorm.Expr("total_favorated + ?", num)).
		Error
	// TotalFavorated
	if favorateErr != nil {
		return favorateErr
	}
	if err != nil {
		return err
	}
	return nil
}

func VideoForFavorite(userID string, videos *[]Video) {
	err := db.Joins("JOIN likes ON likes.video_id = videos.id").
		Where("likes.user_id = ?", userID).
		Find(&videos).Error

	if err != nil {
		panic(err)
	}
}
