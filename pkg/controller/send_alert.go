package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/common"

	errorCommon "github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/error"
)

// follows GCP's pubsub payload structure
type SendAlertRequestBody struct {
	Message struct {
		Data common.EmailAlert `json:"data"`
		ID   string            `json:"id"`
	}
	Subscription string `json:"subscription"`
}

// Simple controller that sends out emails when it receives a message from a
// PubSub topic. Will probably want to make this more robust/swap out technology later on.
func (ctrl Controller) SendAlert(c *gin.Context) {
	const location = "Controller.SendAlert"

	// unmarshal the request body
	var (
		err         error
		requestBody SendAlertRequestBody
	)
	if err = c.BindJSON(&requestBody); err != nil {
		errorCommon.RaiseBadRequestError(c, err, location, "Failed to unmarshal request body")
		return
	}

	var (
		failure = requestBody.Message.Data.ParseFailure
		client  = ctrl.Clients.GetSg()

		from    = mail.NewEmail("NFT Marketplace Staging", "noah.mcgill@merkleroot.co")
		subject = fmt.Sprintf("Max Retries Reached for Failure with Document ID of %s", requestBody.Message.Data.DocID)

		plainTextContent = fmt.Sprintf(`
			The Event Parser Service has failed to parse an event for a tokenId of %s on contract %s after 12 retries. Parsing of this event will no longer 
			be retried, and a record of this failure has been added to the "inactiveParseFailures" collection in Firestore.
		`, failure.TokenID, failure.Contract)

		htmlContent = fmt.Sprintf(`
			<h4>Max Retries Reached for Failure with Document ID of %s</h4>
			<p>
				The Event Parser Service has failed to parse an event for a tokenId of %s on contract %s after 12 retries. 
				Parsing of this event will no longer be retried, and a record of this failure has been added to the "inactiveParseFailures" collection in Firestore.
				<br />
			</p>
			<p>
				To debug:
				<ul>
					<li>Go to Firestore</li>
					<li>Go to the inactiveParseFailures collection</li>
					<li>Find the document with the ID of %s</li>
					<li>Read the "error" field to find out what error occurred</li>
				</ul>
			</p>
			<p>
				This parsing failure will not be retried again.
			</p>
		`, requestBody.Message.Data.DocID, failure.TokenID, failure.Contract, requestBody.Message.Data.DocID)

		// define the recipient email addresses - just leaving these hard-coded for now
		contacts = [3]struct {
			name  string
			email string
		}{
			{
				name:  "Noah McGill",
				email: "noah.mcgill@merkleroot.co",
			},
			{
				name:  "Brandon Brown",
				email: "brandon.brown@merkleroot.co",
			},
			{
				name:  "Brandyn Thibault",
				email: "brandyn.thibault@merkleroot.co",
			},
		}
	)

	for _, contact := range contacts {
		var (
			to      = mail.NewEmail(contact.name, contact.email)
			message = mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		)

		if _, err := client.Send(message); err != nil {
			msg := fmt.Sprintf("Error while sending maxRetry email notification about token %s on contract %s to address %s via SendGrid", failure.TokenID, failure.Contract, contact.email)

			errorCommon.RaiseInternalServerError(c, err, location, msg)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Parsing event failure notifications for token %s on contract %s successfully sent", failure.TokenID, failure.Contract)})
}
