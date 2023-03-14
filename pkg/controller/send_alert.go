package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

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

	//g := ctrl.Clients.GetGmail()

	emailHost := ctrl.Env.SMTP_HOST      ///"smtp.gmail.com"
	emailFrom := ctrl.Env.GMAIL_USER     ///email
	emailPassword := ctrl.Env.GMAIL_PASS ///password
	emailPort := ctrl.Env.SMTP_PORT      ///587

	emailAuth := smtp.PlainAuth("", emailFrom, emailPassword, emailHost)

	//TODO: Replace with logic that returns list of failures for the day
	emailBody := "Test Message"

	year, month, day := time.Now().Date()
	date := fmt.Sprintf("%v/%v/%v", year, month, day)
	emailTo := []string{"matthew.norton@merkleroot.co", "brandyn.thibault@merkleroot.co", "brandon.brown@merkleroot.co", "noah.mcgill@merkleroot.co"}
	subject := "Subject: Daily Failure Report - " + date + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + "\n" + emailBody)
	addr := fmt.Sprintf("%s:%v", emailHost, emailPort)

	err := ctrl.Send(addr, emailAuth, emailFrom, emailTo, msg)

	if err != nil {
		errorCommon.RaiseInternalServerError(c, err, location, "Failed to send daily alert")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Daily alert successfully sent")})
}

func (ctrl Controller) getFailures(c *gin.Context) string {
	var response string

	collection := ctrl.Clients.GetMdb().Database(ctrl.Env.Secrets.DB_NAME).Collection("failures")

	result, _ := collection.Find(c, "test")

	for result.Next(c) {
		var failure Failure
		if err := result.Decode(&failure); err != nil {
			log.Fatal(err)
		}
		response = fmt.Sprintf("%v %v\n", response, result)
	}
	if err := result.Err(); err != nil {
		log.Fatal(err)
	}

	return response
}
