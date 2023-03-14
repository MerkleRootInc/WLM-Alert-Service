package client

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/client"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
)

func InitializeClients(ctx context.Context, err *common.ClientInitErr, env *common.Env) IClients {
	var clients = new(Clients)

	if err == nil {
		return &Clients{}
	}

	// secrets manager
	err.Sm = clients.InitSecretsManagerClient(ctx)

	err.Secrets = common.LoadSecrets(ctx, clients.Sm, env)

	// mongodb
	err.Mdb = clients.InitMongoDBClient(ctx, 20, env.Secrets.MDB_USER, env.Secrets.MDB_PASS)

	// cloud storage
	err.Cs = clients.InitCloudStorageClient(ctx, env.GCP_PROJECT_ID)

	clients.Http = &http.Client{}

	return clients
}

// Creating new structs and methods to avoid making changes to the GoCommon package
// TO DO: Move all below into the common package at a later date
// Note: Once CallContract is added to the IEthClient interface in the common package,
// all of the below can be removed
type Clients struct {
	Cs   client.ICloudStorageClient
	Ps   client.IPubSubClient
	Mdb  client.IMongoDBClient
	Eth  client.IEthClient
	Sm   client.ISecretsManagerClient
	Http client.IHttpClient
	Sg   client.ISendGridClient
}

// Creates a new Ethereum client
func (c *Clients) InitEthClient(ctx context.Context, endpoint string) error {
	eth, err := ethclient.Dial(endpoint)
	if err != nil {
		return err
	}

	c.SetEth(client.EthClient(*eth))

	return nil
}

// Creates a new MongoDB client
func (c *Clients) InitMongoDBClient(ctx context.Context, maxPoolSize int, user, pass string) error {
	var (
		mdb *mongo.Client
		err error

		uri = fmt.Sprintf("mongodb://%s:%s@sample.host:27017/?maxPoolSize=%d&w=majority", user, pass, maxPoolSize)
	)

	mdb, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// check the connection
	err = mdb.Ping(ctx, nil)
	if err != nil {
		return err
	}

	c.Mdb = client.MongoDBClient(*mdb)

	return nil
}

func (c *Clients) CloseMongoDBClient() error {
	if c.Mdb == nil {
		return errors.New("cannot close MongoDB client because there is no active MongoDB client")
	}

	if err := c.Mdb.Disconnect(context.TODO()); err != nil {
		return err
	}

	c.Mdb = nil
	return nil
}

// Creates a new Cloud Storage client
func (c *Clients) InitCloudStorageClient(ctx context.Context, gcpProjectId string) error {
	var (
		cs  *storage.Client
		err error
	)

	cs, err = storage.NewClient(ctx)
	if err != nil {
		return err
	}

	c.Cs = client.CloudStorageClient(*cs)

	return nil
}

// Creates a new PubSub client
func (c *Clients) InitPubSubClient(ctx context.Context, gcpProjectId string) error {
	var (
		ps  *pubsub.Client
		err error
	)

	ps, err = pubsub.NewClient(ctx, gcpProjectId)
	if err != nil {
		return err
	}

	c.Ps = client.PubSubClient(*ps)

	return nil
}

// Creates a new Secrets Manager client
func (c *Clients) InitSecretsManagerClient(ctx context.Context) error {
	var (
		sm  *secretmanager.Client
		err error
	)

	sm, err = secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}

	c.Sm = client.SecretsManagerClient(*sm)

	return nil
}

// Retrieves MongoDB client instance
func (c *Clients) GetMdb() client.IMongoDBClient {
	return c.Mdb
}

// Retrieves Cloud Storage client instance
func (c *Clients) GetCs() client.ICloudStorageClient {
	return c.Cs
}

// Retrieves PubSub client instance
func (c *Clients) GetPs() client.IPubSubClient {
	return c.Ps
}

// Retrieves Ethereum client instance
func (c *Clients) GetEth() client.IEthClient {
	return c.Eth
}

// Retrieves HTTP client instance
func (c *Clients) GetHttp() client.IHttpClient {
	return c.Http
}

// Retrieves Secrets Manager client instance
func (c *Clients) GetSm() client.ISecretsManagerClient {
	return c.Sm
}

// Replaces the ethereum client with a new instance
func (c *Clients) SetEth(eth client.IEthClient) {
	c.Eth = eth
}

func (c *Clients) GetSg() client.ISendGridClient {
	return c.Sg
}
