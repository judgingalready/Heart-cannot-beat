package controller

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb() {
	dsn := "root:Qwer1234!@tcp(localhost:3306)/video?charset=utf8mb4&parseTime=True&loc=Local"
	// Connect Mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Create Table: User
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}
