package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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
		AllowMethods: []string{"GET", "PATCH"},
		AllowHeaders: []string{"Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
	}))

	portal := router.Group("/home")
	{
		portal.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
		portal.GET("/portal", GetPortals)
		portal.POST("/portal", CreatePortal)
		portal.PATCH("/portal", UpdatePortal)
		portal.DELETE("/portal/:id", DeletePortal)
	}

	router.Run(":8888")
}

func GetPortals(c *gin.Context) {
	portals, notFound, err := dbClient.GetAllPortals()
	if err != nil {
		status := http.StatusInternalServerError
		if notFound {
			status = http.StatusNotFound
		}
		loggers.LogError(c, err, status, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": portals,
	})
}

func CreatePortal(c *gin.Context) {
	portal := new(db.Portal)
	if bindErr := c.ShouldBindJSON(portal); bindErr != nil {
		err := fmt.Errorf("failed to bind create portal params; %v", bindErr)
		loggers.LogError(c, err, http.StatusBadRequest, false)
		return
	}

	portal.Name = strings.ToLower(portal.Name)
	id, err := dbClient.InsertPortal(portal)
	if err != nil {
		loggers.LogError(c, err, http.StatusInternalServerError, false)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": id,
	})
}

func UpdatePortal(c *gin.Context) {
	update := new(db.UpdatePortal)
	if bindErr := c.ShouldBindJSON(update); bindErr != nil {
		err := fmt.Errorf("failed to bind update portal params; %v", bindErr)
		loggers.LogError(c, err, http.StatusBadRequest, false)
		return
	}

	updatedPortal, notfound, err := dbClient.ModifyPortal(update)
	if err != nil {
		status := http.StatusBadRequest
		if notfound {
			status = http.StatusNotFound
		}
		loggers.LogError(c, err, status, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedPortal,
	})
}

func DeletePortal(c *gin.Context) {
	portalID := c.Param("id")
	if portalID == "" {
		err := errors.New("missing id param in request")
		loggers.LogError(c, err, http.StatusBadRequest, false)
		return
	}

	err := dbClient.RemovePortal(portalID)
	if err != nil {
		loggers.LogError(c, err, http.StatusInternalServerError, false)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
