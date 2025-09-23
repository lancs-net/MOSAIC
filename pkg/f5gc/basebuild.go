package f5gc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

type Base struct {
	Name        string `json:"name"`
	ImageStatus string `json:"status"`
}

var NFLIST = []string{"smf", "amf", "upf", "udr", "udm", "nrf", "nssf", "n3iwf", "pcf", "ausf", "chf", "n3iwue", "nef", "tngf"}

var HISTORY = []string{}

var BaseImages = []Base{}

func BaseBuild(c *gin.Context) {
	path := c.Param("path")
	var opts docker.ImageBuildOpts
	var webOpts docker.ImageBuildOpts

	d := docker.LocalClient
	opts.Context = path
	opts.Tags = []string{"free5gc/base:latest"}

	webOpts.Context = path
	webOpts.Dockerfile = "Dockerfile.nf.webconsole"
	webOpts.Tags = []string{"free5gc/webconsole-base:latest"}

	err := d.ImageBuild(opts)
	if err != nil {
		fmt.Println("Base Build Failed...")
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Base Build",
			"status":  "failed",
			"error":   err.Error(),
		})
		// return err
	}

	history, err := d.GetImageHistory(opts.Tags[0])
	if err != nil {
		fmt.Println("Image History Fetch Failed...")
		fmt.Println(err)
	}
	HISTORY = append(HISTORY, history...)
	fmt.Println(opts.Tags[0] + " Image History Recorded...")

	// opts = webBase()
	err = d.ImageBuild(webOpts)
	if err != nil {
		fmt.Println("Webconsole Base Build Failed...")
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Web Console Base Build",
			"status":  "failed",
			"error":   err.Error(),
		})
		// return err
	}

	history, err = d.GetImageHistory(opts.Tags[0])
	if err != nil {
		fmt.Println("Image History Fetch Failed...")
		fmt.Println(err)
	}
	HISTORY = append(HISTORY, history...)
	fmt.Println(opts.Tags[0] + " Image History Recorded...")

	for i := range NFLIST {
		opts = nfBase(NFLIST[i], path)
		err = d.ImageBuild(opts)
		if err != nil {
			fmt.Println(NFLIST[i] + " Base Build Failed...")
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": NFLIST[i] + " Base Build",
				"status":  "failed",
				"error":   err.Error(),
			})
			// return err
		}

		history, err := d.GetImageHistory(opts.Tags[0])
		if err != nil {
			fmt.Println("Image History Fetch Failed...")
			fmt.Println(err)
		}
		HISTORY = append(HISTORY, history...)
		fmt.Println(NFLIST[i] + " Image History Recorded...")
	}

	err = d.VolumeCreate("dbdata")
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Volume Build",
			"status":  "failed",
			"error":   err,
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Base Build",
		"status":  "successful",
	})

}
