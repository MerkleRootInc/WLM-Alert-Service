package common

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

import (
	abiCommon "github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/model"
)

/*
Errors
*/

type ClientInitErr struct {
	Cs      error
	Eth     error
	Mdb     error
	Sm      error
	Ps      error
	Secrets error
}

/*
Email Alert Data
*/

type EmailAlert struct {
	DocID        string                 `json:"docId" binding:"required,alphanum"`
	ParseFailure abiCommon.ParseFailure `json:"failure"`
}

// DB model for an Activity (Transfer, Sale or Mint)
type Activity struct {
	ID    *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	To    common.Address      `json:"to,omitempty" bson:"to,omitempty"`
	From  common.Address      `json:"from,omitempty" bson:"from,omitempty"`
	Price *int64              `json:"price,omitempty" bson:"from,omitempty"`

	// Unix format
	Timestamp *int64 `json:"timestamp,omitempty" bson:"timestamp,omitempty"`

	// Activities are stored as subsets to the token doc, so this field does
	// not need to be stored within the activity sub-document.
	TokenID *big.Int `json:"-" bson:"-"`
}

// Payload of a Pub/Sub message, with JSON/BSON bindings and field validation compatible with the Gin framework
// Taken from GCP's example here: https://github.com/GoogleCloudPlatform/golang-samples/blob/main/run/pubsub/main.go#L53-L59
type PubSubMsg struct {
	Message struct {
		ID          string             `json:"id" binding:"required,alphanum"`
		PublishTime time.Time          `json:"publishTime" binding:"required"`
		Attributes  *map[string]string `json:"attributes" binding:"omitempty"`

		// Data as a string represents a base64-encoded marsheling of this package's EthEventData struct
		Data []byte `json:"data" binding:"required"`

		// An automatically-set field representing the number of times a message has been delivered - this might be nil
		DeliveryAttempt *int `json:"deliveryAttempt" binding:"required"`
	} `json:"message" binding:"required"`
	Subscription string `json:"subscription" binding:"required,alphanum"`
}
