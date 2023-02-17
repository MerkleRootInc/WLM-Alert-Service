package main

import (
	"context"
	"fmt"

	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/client"
	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/common"
	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/controller"
	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/middleware"
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
		initClients   = middleware.InitClients(clients, clientsErr)
		bindContracts = middleware.BindContract(clients)
	)

	router.Use(initClients)
	router.Use(bindContracts)

	router.POST("/parse-event", ctrl.ParseEvent)

	return router
}
