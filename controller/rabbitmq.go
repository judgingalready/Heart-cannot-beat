package controller

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FollowPublish(name, message string) {
	ch, err := rabbitmqConnect.Channel()
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	ch.Publish(
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

func FollowConsumer(name string) {
	ch, err := rabbitmqConnect.Channel()
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}

	var forever chan struct{}

	switch name {
	case "FollowAdd": // 关注的消费操作
		go FollowAddConsume(msgs)
	case "FollowDel": // 取关的消费操作
		go FollowDelConsume(msgs)
	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")
	<-forever
}

// 关注功能
func FollowAddConsume(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		toUserId, _ := strconv.ParseInt(params[1], 10, 64)

		CreateFollowErr := db.Create(&Relation{Follow: toUserId, Follower: userId}).Error
		if CreateFollowErr != nil {
			log.Println("Relation数据插入错误：", CreateFollowErr.Error(), userId, toUserId)
		}
	}
	// user.FollowCount += 1
	// userSaveErr := db.Save(&user).Error
	// if userSaveErr != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
	// 	return
	// }

	// toUser.FollowerCount += 1
	// to_userSaveErr := db.Save(&toUser).Error
	// if to_userSaveErr != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
	// 	return
	// }
}

// 取关功能
func FollowDelConsume(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		toUserId, _ := strconv.ParseInt(params[1], 10, 64)

		DeleteFollowErr := db.Where("follow = ? AND follower = ?", toUserId, userId).Delete(&Relation{}).Error
		if DeleteFollowErr != nil {
			log.Println("Relation数据删除错误：", DeleteFollowErr.Error(), userId, toUserId)
		}
	}
	// if user.FollowCount > 0 {
	// 	user.FollowCount -= 1
	// }
	// userSaveErr := db.Save(&user).Error
	// if userSaveErr != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
	// 	return
	// }

	// if to_user.FollowerCount > 0 {
	// 	to_user.FollowerCount -= 1
	// }
	// to_userSaveErr := db.Save(&to_user).Error
	// if to_userSaveErr != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Users数据库更新失败"})
	// 	return
	// }
}
