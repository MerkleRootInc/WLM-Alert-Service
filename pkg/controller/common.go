package controller

import (
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/client"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"net/smtp"
)

/*
Controllers
*/
type Send func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error

type Controller struct {
	Clients client.IClients
	Env     *common.Env
	Send
	Done chan bool
}
