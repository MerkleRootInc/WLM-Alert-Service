package controller

import (
	"errors"
	"github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"os"
	"testing"
)

func TestAlerts(t *testing.T) {
	suite.Run(t, &AlertTestSuite{})
}

func (suite *AlertTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = &http.Request{}

	c := gomock.NewController(suite.T())
	defer c.Finish()
	InitializeClientMocks(&clients, c)

}

func (suite *AlertTestSuite) TestSendAlertSuccess() {
	env := common.Env{
		SMTP_PORT:            os.Getenv("SMTP_PORT"),
		STAGE:                os.Getenv("ENVIRONMENT"),
		GCP_PROJECT_ID:       os.Getenv("PROJECT_ID"),
		SMTP_HOST:            os.Getenv("SMTP_HOST"),
		GMAIL_SECRET_NAME:    "",
		GMAIL_PASS_NAME:      "",
		MDB_USER_SECRET_NAME: "",
		MDB_PASS_SECRET_NAME: "",
		DB_SECRET_NAME:       "",
		Secrets: common.Secrets{
			GMAIL_USER: os.Getenv("GMAIL_USER"),
			GMAIL_PASS: os.Getenv("GMAIL_PASS"),
		},
	}

	ctrl := Controller{
		Clients: clients.Clients,
		Env:     &env,
		Send: func(string, smtp.Auth, string, []string, []byte) error {
			return nil
		},
	}

	ctrl.SendAlert(ctx)
	suite.Equal(http.StatusOK, ctx.Writer.Status(), "incorrect status code returned")
}

func (suite *AlertTestSuite) TestSendAlertFailure() {
	env := common.Env{
		SMTP_PORT:            os.Getenv("SMTP_PORT"),
		STAGE:                os.Getenv("ENVIRONMENT"),
		GCP_PROJECT_ID:       os.Getenv("PROJECT_ID"),
		SMTP_HOST:            os.Getenv("SMTP_HOST"),
		GMAIL_SECRET_NAME:    "",
		GMAIL_PASS_NAME:      "",
		MDB_USER_SECRET_NAME: "",
		MDB_PASS_SECRET_NAME: "",
		DB_SECRET_NAME:       "",
		Secrets: common.Secrets{
			GMAIL_USER: os.Getenv("GMAIL_USER"),
			GMAIL_PASS: os.Getenv("GMAIL_PASS"),
		},
	}

	ctrl := Controller{
		Clients: clients.Clients,
		Env:     &env,
		Send: func(string, smtp.Auth, string, []string, []byte) error {
			return errors.New("test error")
		},
	}

	ctrl.SendAlert(ctx)
	suite.Equal(http.StatusInternalServerError, ctx.Writer.Status(), "incorrect status code returned")
}
