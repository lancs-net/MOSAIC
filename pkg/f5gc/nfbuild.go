package f5gc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/common"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

func NfBuild(c *gin.Context) {
	d := docker.LocalClient
	var nfjson common.NFJSON

	if err := c.BindJSON(&nfjson); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload",
			"status":  "failed",
			"error":   err.Error(),
		})
	}

	nfinfo := nfjson.NFInfos
	path := nfjson.Path

	conOpts := dbConf()
	imgOpts := docker.ImageBuildOpts{}
	err := BuildStart(d, imgOpts, conOpts)
	if err != nil {
		fmt.Println("DB Build Failed...")
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "DB Build",
			"status":  "failed",
			"error":   err.Error(),
		})
	}
	fmt.Println("DB Build Successful...")

	imgOpts, conOpts = webuiConf(path)
	err = BuildStart(d, imgOpts, conOpts)
	if err != nil {
		fmt.Println("WebUI Build Failed...")
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "WebUI Build",
			"status":  "failed",
			"error":   err.Error(),
		})
	}
	fmt.Println("WebUI Build Successful...")

	for i := range nfinfo {
		for j := 1; j <= nfinfo[i].NumInstances; j++ {
			imgOpts, conOpts := GetConf(nfinfo[i], path, j)
			err := BuildStart(d, imgOpts, conOpts)
			if err != nil {
				fmt.Println(nfinfo[i].Name + " Build Failed...")
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"message": nfinfo[i].Name + " NF Build",
					"status":  "failed",
					"error":   err.Error(),
				})
			}
		}
		fmt.Println(nfinfo[i].Name + " Build Successful...")
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "NF Build",
		"status":  "successful",
	})
}

func BuildStart(d *docker.DockerClient, imgOpts docker.ImageBuildOpts, conOpts docker.ContainerCreateOpts) error {
	if conOpts.Name == "mongodb" {
		err := d.ImagePull("docker.io/library/mongo")
		if err != nil {
			fmt.Println("DB Image Build Failed...")
			fmt.Println(err)
			return err
		}
	} else {
		err := d.ImageBuild(imgOpts)
		if err != nil {
			fmt.Println(conOpts.Name + " Image Build Failed...")
			fmt.Println(err)
			return err
		}
	}

	fmt.Println(conOpts.Name + " Image Build Successful...")
	conID, err := d.ContainerCreate(conOpts)
	if err != nil {
		fmt.Println(conOpts.Name + " Container Create Failed...")
		fmt.Println(err)
		return err
	}
	fmt.Println(conOpts.Name + " Container Create Successful...")

	err = d.ContainerStart(conID)
	if err != nil {
		fmt.Println(conOpts.Name + " Container Start Failed...")
		fmt.Println(err)
		return err
	}
	fmt.Println(conOpts.Name + " Container Start Successful...")

	return nil
}
