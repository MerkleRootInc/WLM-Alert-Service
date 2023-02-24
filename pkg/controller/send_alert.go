package controller

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"net/http"

	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/gin-gonic/gin"

	errorCommon "github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/error"
)

// SendAlertRequestBody follows GCP's pubsub payload structure
type SendAlertRequestBody struct {
	Message struct {
		Data common.EmailAlert `json:"data"`
		ID   string            `json:"id"`
	}
	Subscription string `json:"subscription"`
}

// SendAlert Simple controller that sends out emails when it receives a message from a
// PubSub topic. Will probably want to make this more robust/swap out technology later on.
func (ctrl Controller) SendAlert(c *gin.Context) {
	const location = "Controller.SendAlert"

	g := ctrl.Clients.GetGmail()

	emailBody := "Test Message"

	var alert gmail.Message

	emailTo := "To: menorton15@hotmail.com\r\n"
	subject := "Subject: " + "Test Email form Gmail API using OAuth" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	alert.Raw = base64.URLEncoding.EncodeToString(msg)
	g.Send(&alert)
	// unmarshal the request body
	var (
		err         error
		requestBody SendAlertRequestBody
	)
	if err = c.BindJSON(&requestBody); err != nil {
		errorCommon.RaiseBadRequestError(c, err, location, "Failed to unmarshal request body")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Daily alert successfully sent")})
}
