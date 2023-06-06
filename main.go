package main

import (
	"github.com/gin-gonic/gin"
	db "github.com/mohamedmuhsinJ/loginAssignment/Db"
	"github.com/mohamedmuhsinJ/loginAssignment/controllers"
)

func init() {
	db.ConnectToDb()
	controllers.SyncDb()
}

func main() {
	router := gin.Default()

	router.POST("/register", controllers.Register)
	router.GET("/", controllers.Home)
	router.Run(":8080")

}
