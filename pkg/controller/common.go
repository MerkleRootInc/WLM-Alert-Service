package controller

import (
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/client"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
)

/*
Controllers
*/

type Controller struct {
	Clients client.IClients
	Env     common.Env
	Done    chan bool
}
