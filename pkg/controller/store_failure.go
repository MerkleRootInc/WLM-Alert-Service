package controller

import (
	"fmt"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/client"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/ethereum/go-ethereum/core/types"
	"net/http"
	"time"

	errorCommon "github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/error"
	"github.com/gin-gonic/gin"
)

var CurrentTime = func() time.Time {
	return time.Now()
}

// StoreFailureRequestBody follows GCP's pubsub payload structure
type StoreFailureRequestBody struct {
	Message struct {
		Data        types.Log `json:"data"`
		MessageID   string    `json:"id"`
		PublishTime time.Time
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type Failure struct {
	Timestamp time.Time `bson:"timestamp"`
	Log       types.Log `bson:"log"`       // holds all the event data, including tx hash, log index, and contract address
	MessageID string    `bson:"messageID"` // pub/sub message ID

	// below field values come from the Log (pulling them out into separate fields for ease of viewing)
	Tx       string `bson:"txn"`
	LogIndex string `bson:"logIndex"`
	Contract string `bson:"contract"`
}

// StoreFailure Simple controller that stores failures to mongoDB.
func (ctrl Controller) StoreFailure(c *gin.Context) {
	const (
		location   = "Controller.StoreFailure"
		collection = "failures"
	)

	var (
		err         error
		requestBody StoreFailureRequestBody
		mongo       client.IMongoDBClient
		failures    client.IMongoDBCollection
	)

	var database = common.DB_NAME

	if err = c.BindJSON(&requestBody); err != nil {
		errorCommon.RaiseBadRequestError(c, err, location, "Failed to unmarshal request body")
	}

	var message = requestBody.Message

	mongo = ctrl.Clients.GetMdb()

	failures = mongo.Database(database).Collection(collection)

	newFailure := parseFailure(message.Data, message.MessageID)

	if _, err = failures.InsertOne(c, newFailure); err != nil {
		errorCommon.RaiseInternalServerError(c, err, location, "Failed to store parse event failure")
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Parsing event failure for txn %s on contract %s successfully stored", newFailure.Tx, newFailure.Contract)})

}

func parseFailure(data types.Log, messageID string) *Failure {
	failure := Failure{
		Timestamp: CurrentTime(),
		Log:       data,
		MessageID: messageID,
		Tx:        data.TxHash.String(),
		LogIndex:  fmt.Sprint(data.Index),
		Contract:  data.Address.String(),
	}
	return &failure
}

func (p *Failure) Matches(x interface{}) bool {
	///reflectedValue := reflect.DeepEqual(reflect.ValueOf(x), p)

	return true
}
