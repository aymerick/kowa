package server

import (
	"github.com/aymerick/kowa/models"

	"github.com/gin-gonic/gin"
)

// endpoint: list all sites
func listSites(c *gin.Context) {
	c.JSON(200, gin.H{"sites": models.AllSites()})
}
