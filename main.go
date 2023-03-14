package main

import (
	"context"
	"fmt"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/client"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/controller"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"net/smtp"
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
		env        = common.LoadEnv()
		clientsErr = common.ClientInitErr{}
	)

	if env.STAGE == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	clients := client.InitializeClients(ctx, &clientsErr, env)
	ctrl := controller.Controller{Clients: clients, Env: env, Send: smtp.SendMail}

	// set up the middleware
	var (
		initClients = middleware.InitClients(clients, clientsErr)
	)

	router.Use(initClients)

	router.POST("/send-alert", ctrl.SendAlert)

	return router
}
