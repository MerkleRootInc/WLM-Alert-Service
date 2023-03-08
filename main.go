package main

import (
	"context"
	"fmt"

	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/client"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/controller"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize the gin router
	router := setupRouter()

	// run the server on port 8080 by default
	router.Run(fmt.Sprintf(":%s", common.PORT))
}

func setupRouter() *gin.Engine {
	// create router and controller instance with dependencies
	var (
		router     = gin.Default()
		ctx        = context.Background()
		clientsErr = common.ClientInitErr{}
		clients    = client.InitializeClients(ctx, &clientsErr)
		ctrl       = controller.Controller{Clients: clients}
	)

	// set up the middleware
	var (
		initClients = middleware.InitClients(clients, clientsErr)
	)

	router.Use(initClients)

	router.POST("/send-alert", ctrl.SendAlert)

	return router
}
