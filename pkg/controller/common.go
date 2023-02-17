package controller

import (
	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/client"
)

/*
Controllers
*/

type Controller struct {
	Clients client.IClients
	Done    chan bool
}
