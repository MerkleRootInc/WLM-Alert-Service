package controller

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/test"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"time"
)

var privateKey, err = crypto.GenerateKey()
var publicKey = privateKey.Public()
var publicKeyECDSA = publicKey.(*ecdsa.PublicKey)
var address = crypto.PubkeyToAddress(*publicKeyECDSA)

var TestLog1 = types.Log{
	Address:     address,
	Topics:      []common.Hash{common.HexToHash("0")},
	Data:        []byte{},
	BlockNumber: uint64(0),
	TxHash:      common.HexToHash("0"),
	TxIndex:     uint(0),
	BlockHash:   common.HexToHash("0"),
	Index:       uint(0),
	Removed:     false,
}

func GetTestRequest(log *types.Log) ([]byte, error) {
	return json.Marshal(StoreFailureRequestBody{
		Message: struct {
			Data        types.Log `json:"data"`
			MessageID   string    `json:"id"`
			PublishTime time.Time
		}{
			Data:        TestLog1,
			MessageID:   "test",
			PublishTime: time.Now(),
		},
		Subscription: "testSubscription",
	})
}

func InitializeClientMocks(mock *test.ClientMock, c *gomock.Controller) {
	*mock = test.ClientMock{
		Clients: test.NewMockIClients(c),
		Abis:    test.AbisMock{},
		Auth:    test.AuthMock{},
		Fs:      test.FsMock{},
		Cs:      test.CsMock{},
		Eth:     test.EthMock{},
		Http:    test.HttpMock{},
		Mdb:     test.MdbMock{Client: test.NewMockIMongoDBClient(c), Database: test.NewMockIMongoDBDatabase(c), Collection: test.NewMockIMongoDBCollection(c)},
		Ps:      test.PsMock{},
		Redis:   test.RedisMock{},
		Sm:      test.SmMock{},
		Sg:      test.SgMock{},
	}
}
