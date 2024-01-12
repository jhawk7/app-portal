package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jhawk7/app-portal/internal/pkg/db"
	"github.com/jhawk7/app-portal/internal/pkg/loggers"
)

var dbClient *db.DBClient

func initializeDB() {
	client, dbErr := db.InitDB()
	loggers.LogError(nil, dbErr, 0, true)
	dbClient = client
}

func main() {
	initializeDB()
	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
	}))

	portal := router.Group("/home")
	{
		portal.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
		portal.GET("/portals", GetPortals)
		portal.POST("/portal", CreatePortal)
		portal.PATCH("/portal/:id", UpdatePortal)
		portal.DELETE("/portal/:id", DeletePortal)
	}

	router.Run(":8888")
}

func GetPortals(c *gin.Context) {

}

func CreatePortal(c *gin.Context) {

}

func UpdatePortal(c *gin.Context) {

}

func DeletePortal(c *gin.Context) {

}
