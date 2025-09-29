package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lancs-net/MOSAIC/pkg/common"
	"github.com/lancs-net/MOSAIC/pkg/f5gc"
)

func main() {
	requestRouter := gin.Default()

	requestRouter.GET("/", common.Home)

	// For builidng and running the networks for the deployment
	// needs a JSON input in the form {"net_name":"subnet", ...}
	requestRouter.POST("/net", common.NetworkInit)
	requestRouter.DELETE("/net/:network", common.NetworkRemove)

	// For building the base images for all network functions
	requestRouter.POST("/f5gc/base/:path", f5gc.BaseBuild)
	requestRouter.DELETE("/f5gc/base", f5gc.RemoveBase)

	// For building and running the network functions
	// needs a JSON input in the form {"nf_name":{"number": 1, "network": "" , ...}}
	requestRouter.POST("/f5gc", f5gc.NfBuild)
	requestRouter.PUT("/f5gc/:step/:nf", f5gc.Nf)
	requestRouter.DELETE("/f5gc/:nf", f5gc.NfRemove)

	// To get the status of a network function
	// needs JSON used for deployment
	requestRouter.GET("/status", common.NfStatus)

	requestRouter.Run("0.0.0.0:8000")
}
