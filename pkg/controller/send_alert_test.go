package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/googleapi"
	"net/http"
	"net/http/httptest"
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

	ctrl := Controller{Clients: clients.Clients}

	clients.Clients.EXPECT().GetGmail().Return(clients.Gmail.Client)
	clients.Gmail.Client.EXPECT().Send(gomock.Any()).Return(&googleapi.ServerResponse{HTTPStatusCode: http.StatusOK, Header: http.Header{}}, nil)

	ctrl.SendAlert(ctx)
	suite.Equal(http.StatusOK, ctx.Writer.Status(), "incorrect status code returned")
}

func (suite *AlertTestSuite) TestSendAlertFailure() {

	ctrl := Controller{Clients: clients.Clients}

	clients.Clients.EXPECT().GetGmail().Return(clients.Gmail.Client)
	clients.Gmail.Client.EXPECT().Send(gomock.Any()).Return(&googleapi.ServerResponse{
		HTTPStatusCode: http.StatusInternalServerError, Header: http.Header{}}, errors.New("test error message"))

	ctrl.SendAlert(ctx)
	suite.Equal(http.StatusInternalServerError, ctx.Writer.Status(), "incorrect status code returned")
}
