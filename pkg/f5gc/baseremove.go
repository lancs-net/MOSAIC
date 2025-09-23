package f5gc

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

func RemoveBase(c *gin.Context) {
	d := docker.LocalClient
	for i := range NFLIST {
		err := Remove(d, NFLIST[i])
		if err != nil {
			fmt.Println(NFLIST[i] + " Remove Failed...")
			fmt.Println(err)
			c.IndentedJSON(
				http.StatusInternalServerError, gin.H{
					"message": NFLIST[i] + " Base Removal",
					"status":  "failed",
					"error":   err.Error(),
				})
		} else {
			fmt.Println(NFLIST[i] + " Remove Successful...")
			c.IndentedJSON(
				http.StatusInternalServerError, gin.H{
					"message": "NF Base Removal",
					"status":  "Successful",
				})
		}
	}

	err := Remove(d, "webconsole")
	if err != nil {
		fmt.Println("Webconsole Remove Failed...")
		fmt.Println(err)
		c.IndentedJSON(
			http.StatusInternalServerError, gin.H{
				"message": "Webconsole Base Removal",
				"status":  "failed",
				"error":   err.Error(),
			})
	} else {
		fmt.Println("Webconsole Remove Successful...")
		c.IndentedJSON(
			http.StatusInternalServerError, gin.H{
				"message": "Webconsole Base Removal",
				"status":  "Successful",
			})
	}

	slices.Reverse(HISTORY)
	fmt.Println("Removing Intermediate Images...")
	for _, id := range HISTORY {
		err = d.ImageRemove(id)
		if err != nil {
			fmt.Println("Intermediate Image Removal Failed...")
			fmt.Println(err)
		}
	}

	err = Remove(d, "base")
	if err != nil {
		fmt.Println("Base Removal Failed...")
		fmt.Println(err)
	} else {
		fmt.Println("Base Removal Successful...")
	}

	err = d.VolumeRemove("dbdata")
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Volume Removal",
			"status":  "failed",
			"error":   err,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Base Removal",
		"status":  "successful",
	})

}
