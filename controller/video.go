package controller

import "fmt"

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
	fmt.Println(videos)
}
