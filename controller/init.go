package controller

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB //全局的db对象

func InitDb() {
	//配置MySQL连接参数
	username := "root"      //账号
	password := "Qwer1234!" //密码
	host := "127.0.0.1"     //数据库地址，可以是Ip或者域名
	port := 3306            //数据库端口
	Dbname := "DOUSHENG"    //数据库名
	timeout := "10s"        //连接超时，10

	//拼接dsn参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	// Connect Mysql
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Create Table: User
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Video{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Account{})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()

	//设置数据库连接池参数    sqlDB.SetMaxOpenConns(100)   //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
}

/*
*获取gorm db对象，其他包需要执行数据库查询的时候，只要通过tools.getDB()获取db对象即可。不用担心协程并发使
*用同样的db对象会共用同一个连接，db对象在调用他的方法的时候会从数据库连接池中获取新的连接。
 */
func GetDB() *gorm.DB {
	return db
}
