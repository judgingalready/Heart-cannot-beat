package main

import (
	"simple-demo/controller"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	go service.RunMessageServer()
	controller.InitDb()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
