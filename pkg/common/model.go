package common

import (
	"fmt"
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
	Cs  error
	Eth error
	Mdb error
	Sm  error
	Ps  error
}

/*
Email Alert Data
*/

type EmailAlert struct {
	DocID        string                 `json:"docId" binding:"required,alphanum"`
	ParseFailure abiCommon.ParseFailure `json:"failure"`
}

/*
Contract Bindings
*/

// Creating new interfaces to avoid making changes to the GoCommon package
// TO DO: Move these into the common package at a later date

type ContractsClient struct {
	EIP165  IEIP165Contract
	ERC721  IERC721Contract
	ERC1155 IERC1155Contract
}

// Retrieves the contract as an EIP165 instance
func (c *ContractsClient) GetEIP165Contract() IEIP165Contract {
	return c.EIP165
}

// Retrieves the contract as an EIP165 instance
func (c *ContractsClient) GetERC721Contract() IERC721Contract {
	return c.ERC721
}

// Retrieves the contract as an EIP165 instance
func (c *ContractsClient) GetERC1155Contract() IERC1155Contract {
	return c.ERC1155
}

/*
Events
*/

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

// Creating new structs to avoid making changes to the GoCommon package
// TO DO: Move to/integrate with full structs in the common package at a later date

/*
Tokens
*/

type TokenDocument struct {
	ID             *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Contract       string              `json:"contract,omitempty" bson:"contract,omitempty"`
	TokenID        string              `json:"tokenId,omitempty" bson:"tokenId,omitempty"`
	Likes          *int                `json:"likes,omitempty" bson:"likes,omitempty"`
	Views          *int                `json:"views,omitempty" bson:"views,omitempty"`
	IsListed       *bool               `json:"isListed,omitempty" bson:"isListed,omitempty"`
	Featured       *bool               `json:"featured,omitempty" bson:"featured,omitempty"`
	ActiveOffers   *[]Offer            `json:"activeOffers,omitempty" bson:"activeOffers,omitempty"`
	CurrentAuction *Auction            `json:"currentAuction,omitempty" bson:"currentAuction,omitempty"`
	History        *TokenHistory       `json:"history,omitempty" bson:"history,omitempty"`

	// The below fields represent the metadata we will store for a token. Adhere's to OpenSea's accepted
	// metadata standards, which includes the Enjin metadata format. JSON field names mirror URI field
	// names, and BSON field names represent the names under which we'll store the data fields in the
	// database (to maintain MongoDB's camel case standard).
	// https://docs.opensea.io/docs/metadata-standards
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1155.md#erc-1155-metadata-uri-json-schema
	AnimationURL    string `json:"animation_url,omitempty" bson:"animationUrl,omitempty"`
	Attributes      string `json:"attributes,omitempty" bson:"attributes,omitempty"` // (represented as array of objs)
	BackgroundColor string `json:"background_color,omitempty" bson:"backgroundColor,omitempty"`
	Color           string `json:"color,omitempty" bson:"color,omitempty"`       // Enjin field
	Decimals        *int   `json:"decimals,omitempty" bson:"decimals,omitempty"` // Enjin field
	Description     string `json:"description,omitempty" bson:"description,omitempty"`
	ExternalURL     string `json:"external_url,omitempty" bson:"externalUrl,omitempty"`
	Image           string `json:"image,omitempty" bson:"image,omitempty"`
	CompressedImage string `json:"compressedImage,omitempty" bson:"compressedImage,omitempty"`
	ImageData       string `json:"image_data,omitempty" bson:"imageData,omitempty"`
	Localization    string `json:"localization,omitempty" bson:"localization,omitempty"` // Enjin field (represented as obj)
	Name            string `json:"name,omitempty" bson:"name,omitempty"`
	Properties      string `json:"properties,omitempty" bson:"properties,omitempty"` // Enjin field (represented as obj)
	YouTubeURL      string `json:"youtube_url,omitempty" bson:"youtubeUrl,omitempty"`
}

func (t TokenDocument) GetTokenID() string {
	return t.ID.String()
}

func (t TokenDocument) ToString() string {
	return fmt.Sprintf("%+v", t)
}

type TokenHistory struct {
	Activities     *[]Activity `json:"activities,omitempty" bson:"activities,omitempty"`
	InactiveOffers *[]Offer    `json:"inactiveOffers,omitempty" bson:"inactiveOffers,omitempty"`
	Auctions       *[]Auction  `json:"auctions,omitempty" bson:"auctions,omitempty"`
}

type Offer struct {
	OrderNonce  string     `json:"orderNonce,omitempty" bson:"orderNonce,omitempty"`
	Maker       string     `json:"maker,omitempty" bson:"maker,omitempty"`
	Amount      *float64   `json:"amount,omitempty" bson:"amount,omitempty"`
	Accepted    *bool      `json:"accepted,omitempty" bson:"accepted,omitempty"`
	Date        *time.Time `json:"date,omitempty" bson:"date,omitempty"`
	OfferStart  *time.Time `json:"offerStart,omitempty" bson:"offerStart,omitempty"`
	OfferLength *Length    `json:"offerLength,omitempty" bson:"offerLength,omitempty"`
}

type Auction struct {
	AuctionID     string     `json:"auctionId,omitempty" bson:"auctionId,omitempty"`
	IsActive      *bool      `json:"isActive,omitempty" bson:"isActive,omitempty"`
	NumOfBids     *int       `json:"numOfBids,omitempty" bson:"numOfBids,omitempty"`
	HighBid       *Bid       `json:"highBid,omitempty" bson:"highBid,omitempty"`
	AuctionStart  *time.Time `json:"auctionStart,omitempty" bson:"auctionStart,omitempty"`
	AuctionLength *Length    `json:"auctionLength,omitempty" bson:"auctionLength,omitempty"`
}

type Bid struct {
	OrderNonce string     `json:"orderNonce,omitempty" bson:"orderNonce,omitempty"`
	AuctionID  string     `json:"auctionId,omitempty" bson:"auctionId,omitempty"`
	Maker      string     `json:"maker,omitempty" bson:"maker,omitempty"`
	Amount     *float64   `json:"amount,omitempty" bson:"amount,omitempty"`
	Accepted   *bool      `json:"accepted,omitempty" bson:"accepted,omitempty"`
	Date       *time.Time `json:"date,omitempty" bson:"date,omitempty"`
}

type Length struct {
	Value *int   `json:"value,omitempty" bson:"value,omitempty"`
	Unit  string `json:"unit,omitempty" bson:"unit,omitempty"`
}

// This differs from the `TokenHistory` struct as it's designed to be a stand-alone document
// in a collection rather than an embedded document.
type TokenHistoryDocument struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Contract string `json:"contract,omitempty" bson:"contract,omitempty"`
	TokenID  string `json:"tokenId,omitempty" bson:"tokenId,omitempty"`

	Activities     *[]Activity `json:"activities,omitempty" bson:"activities,omitempty"`
	InactiveOffers *[]Offer    `json:"inactiveOffers,omitempty" bson:"inactiveOffers,omitempty"`
	Auctions       *[]Auction  `json:"auctions,omitempty" bson:"auctions,omitempty"`
}

/*
Misc
*/

type MediaType string

const (
	IMAGE     MediaType = "IMAGE"
	IMAGE_GIF MediaType = "IMAGE_GIF"
	VIDEO     MediaType = "VIDEO"
	UNKNOWN   MediaType = "UNKNOWN"
)

type CustomParserConfig struct {
	ID        string `json:"id" bson:"_id"`
	Contract  string `json:"contract" bson:"contract"`
	TopicID   string `json:"topicId" bson:"topicId"`
	ParserUrl string `json:"parserUrl" bson:"parserUrl"`
}
