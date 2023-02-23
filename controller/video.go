package controller

import "gorm.io/gorm"

const videoCount = 30

func SearchVideoForFeed(videos *[]Video) {
	err := db.Preload("Author").Limit(videoCount).Find(videos).Error
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

func VideoForAction(video_id int64, videos *[]Video, num int32) {
	err := db.Where("id = ?", video_id).Find(videos).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", num)).Error
	if err != nil {
		panic(err)
	}
}

func VideoForFavorite(userID string, videos *[]Video) {
	err := db.Joins("JOIN likes ON likes.video_id = videos.id").
		Where("likes.user_id = ?", userID).
		Find(&videos).Error

	if err != nil {
		panic(err)
	}
}
