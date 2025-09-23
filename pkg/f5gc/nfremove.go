package f5gc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/common"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

func NfRemove(c *gin.Context) {
	d := docker.LocalClient
	nf := c.Param("nf")

	var nfjson common.NFJSON

	if err := c.BindJSON(&nfjson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload",
			"status":  "failed",
			"error":   err.Error(),
		})
	}

	if nf == "all" {
		funcExists, err := RemoveAllCon(d, nfjson.NFInfos)
		if err != nil && funcExists {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "NF Container Removal",
				"status":  "failed",
				"error":   err.Error(),
			})
			return
		} else if err != nil && !funcExists {
			fmt.Println(err)
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "NF Not Found",
				"status":  "failed",
			})
			return
		}
	} else {
		funcExists, err := RemoveCon(d, nf)
		if err != nil && funcExists {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": nf + " Container Removal",
				"status":  "failed",
				"error":   err.Error(),
			})
			return
		} else if err != nil && !funcExists {
			fmt.Println(err)
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": nf + " Not Found",
				"status":  "failed",
			})
			return
		}
	}
}
