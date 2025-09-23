package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NFInfo struct {
	Name         string `json:"name"`
	NumInstances int    `json:"number"`
	Network      string `json:"network"`
}

type NFJSON struct {
	Path    string   `json:"Path"`
	NFInfos []NFInfo `json:"NetworkFunctions"`
}

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to MOSAIC API",
	})
}
