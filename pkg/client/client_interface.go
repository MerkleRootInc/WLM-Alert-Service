package client

import (
	"context"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/client"
)

// Creating new interfaces to avoid making changes to the GoCommon package
// TO DO: Move these into the common package at a later date

type IClients interface {
	InitEthClient(ctx context.Context, endpoint string) error
	InitMongoDBClient(ctx context.Context, maxPoolSize int, user, pass string) error
	InitCloudStorageClient(ctx context.Context, gcpProjectId string) error
	InitPubSubClient(ctx context.Context, gcpProjectId string) error
	InitSecretsManagerClient(ctx context.Context) error

	GetMdb() client.IMongoDBClient
	GetCs() client.ICloudStorageClient
	GetPs() client.IPubSubClient
	GetEth() client.IEthClient
	GetSm() client.ISecretsManagerClient
	GetHttp() client.IHttpClient
	GetSg() client.ISendGridClient

	SetEth(eth client.IEthClient)
	CloseMongoDBClient() error
}
