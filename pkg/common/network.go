package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

type Networks map[string]string

func NetworkInit(c *gin.Context) {
	d := docker.LocalClient
	var net Networks

	if err := c.BindJSON(&net); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for name, subnet := range net {
		var opts docker.NetworkCreateOpts
		opts.Name = name
		opts.Subnet = subnet

		netexists, err := d.NetworkExists(opts.Name)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": opts.Name + " Network Existence Check",
				"status":  "failed",
				"error":   err.Error(),
			})
			return
		}
		if !netexists {
			_, err = d.NetworkCreate(opts)
			if err != nil {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": opts.Name + " Network Create",
					"status":  "failed",
					"error":   err.Error(),
				})
				return
			}
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Networks Initialized", "status": "success"})
}

func NetworkRemove(c *gin.Context) {
	d := docker.LocalClient
	var net Networks

	if err := c.BindJSON(&net); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	network := c.Param("network")

	if network == "" {
		for name, _ := range net {
			netexists, err := d.NetworkExists(name)
			if err != nil {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": name + " Network Existence Check",
					"status":  "failed",
					"error":   err.Error(),
				})
				return
			}
			if netexists {
				err = d.NetworkRemove(name)
				if err != nil {
					fmt.Println(err)
					c.IndentedJSON(http.StatusInternalServerError, gin.H{
						"message": name + " Network Remove",
						"status":  "failed",
						"error":   err.Error(),
					})
					return
				}
			}
		}
	} else {
		netexists, err := d.NetworkExists(network)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": network + " Network Existence Check",
				"status":  "failed",
				"error":   err.Error(),
			})
			return
		}
		if netexists {
			err = d.NetworkRemove(network)
			if err != nil {
				fmt.Println(err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": network + " Network Remove",
					"status":  "failed",
					"error":   err.Error(),
				})
				return
			}
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Networks Removed", "status": "success"})
}
