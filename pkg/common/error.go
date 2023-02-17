package common

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Creating new interfaces to avoid making changes to the GoCommon package
// TO DO: Move these into the common package at a later date

type ErrorMsg string

type AlertErrCode int

type EventParserErrOpts struct {
	// HTTP code that will be appended to response
	Code *AlertErrCode

	// The function name wherein the error occurred
	Location string

	// The message ID of the pub/sub message that failed to be parsed resulting in the error
	MessageID string

	// The contract address corresponding to the event that failed to be parsed resulting in the error
	Address string

	// The tokenId corresponding to the event that failed to be parsed resulting in the error
	TxHash string
}

var (
	StatusInternalServerError AlertErrCode = http.StatusInternalServerError
	StatusBadRequest          AlertErrCode = http.StatusBadRequest
	StatusForbidden           AlertErrCode = http.StatusForbidden
	StatusUnauthorized        AlertErrCode = http.StatusUnauthorized
	StatusNotFound            AlertErrCode = http.StatusNotFound

	errMsgs = map[AlertErrCode]string{
		StatusInternalServerError: "Internal server error.",
		StatusBadRequest:          "Bad request.",
		StatusForbidden:           "Invalid credentials.",
		StatusUnauthorized:        "Unauthorized to perform action.",
		StatusNotFound:            "Not found.",
	}
)

// RaiseAlertErr Appends a generic error message and status code to a Gin response and logs the actual error
func RaiseAlertErr(requestContext *gin.Context, opts EventParserErrOpts, err error) {
	logLine := fmt.Sprintf(": %s", err.Error())

	if opts.TxHash != "" {
		logLine = fmt.Sprintf("[txHash=%s]%s", opts.TxHash, logLine)
	}
	if opts.Address != "" {
		logLine = fmt.Sprintf("[address=%s]%s", opts.Address, logLine)
	}
	if opts.MessageID != "" {
		logLine = fmt.Sprintf("[messageId=%s]%s", opts.MessageID, logLine)
	}
	if opts.Location != "" {
		logLine = fmt.Sprintf("[location=%s]%s", opts.Location, logLine)
	}

	logLine = fmt.Sprintf("Event Parser Service %s", logLine)

	// log the error
	log.Println(logLine)

	// lookup the corresponding generic error message
	var errMsg string
	if opts.Code != nil {
		errMsg = errMsgs[*opts.Code]
		*opts.Code = StatusInternalServerError
	} else {
		errMsg = "Unknown error."
	}

	// append the generic error to the Gin response
	requestContext.JSON(int(*opts.Code), gin.H{"message": errMsg})
}
