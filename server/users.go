package server

import (
	"github.com/aymerick/kowa/models"

	"github.com/gin-gonic/gin"
)

// endpoint: list all users
func listUsers(c *gin.Context) {
	c.JSON(200, gin.H{"users": models.AllUsers()})
}
