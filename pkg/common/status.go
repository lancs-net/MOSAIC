package common

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

type NetworkFunctionStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Image  bool   `json:"image"`
}

func NfStatus(c *gin.Context) {
	d := docker.LocalClient
	var nfjson NFJSON

	if err := c.BindJSON(&nfjson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload",
			"status":  "failed",
			"error":   err.Error(),
		})
	}

	nfinfo := nfjson.NFInfos

	for i := range nfinfo {
		for j := 1; j <= nfinfo[i].NumInstances; j++ {
			status := Status(d, nfinfo[i].Name+strconv.Itoa(j))
			c.IndentedJSON(http.StatusOK, status)
		}
	}
}

func Status(d *docker.DockerClient, function string) NetworkFunctionStatus {
	var status NetworkFunctionStatus
	status.Name = function
	val, err := d.ContainerStatus(function)
	if err == nil && val != "" {
		status.Status = val
		status.Image = true
	} else if err == nil && val == "" {
		val, err := d.ImageExists(function + ":latest")
		if err == nil && val {
			status.Status = "not created"
			status.Image = true
		} else if err == nil && !val {
			status.Status = "not created"
			status.Image = false
		}
	}
	return status
}
