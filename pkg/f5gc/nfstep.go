package f5gc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/common"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

func Nf(c *gin.Context) {
	d := docker.LocalClient
	var nfjson common.NFJSON

	if err := c.BindJSON(&nfjson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload",
			"status":  "failed",
			"error":   err.Error(),
		})
	}

	step := c.Param("step")
	nf := c.Param("nf")

	if step == "stop" {
		if nf == "all" {
			funcExists, err := StopAll(d, nfjson.NFInfos)
			if err != nil && funcExists {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": "NF Stop",
					"status":  "failed",
					"error":   err,
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
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "NF Stop",
				"status":  "successful",
			})
		} else {
			funcExists, err := Stop(d, nf)
			if err != nil && funcExists {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": "NF Stop",
					"status":  "failed",
					"error":   err,
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
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": nf + " Stop",
				"status":  "successful",
			})
		}
	} else if step == "start" {
		if nf == "all" {
			funcExists, err := StartAll(d, nfjson.NFInfos)
			if err != nil && funcExists {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": "NF Start",
					"status":  "failed",
					"error":   err,
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
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "NF Start",
				"status":  "successful",
			})
		} else {
			funcExists, err := Start(d, nf)
			if err != nil && funcExists {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": "NF Start",
					"status":  "failed",
					"error":   err,
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
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": nf + " Start",
				"status":  "successful",
			})
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid step",
			"status":  "failed",
		})
	}
}
