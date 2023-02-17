package controller

import (
	"bytes"
	"errors"
	"fmt"
	abiCommon "github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/model"
	"github.com/MerkleRootInc/NFT-Marketplace-GoCommon/pkg/test"
	parseCommon "github.com/MerkleRootInc/WLM-Alert-Service/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var emailAlert = parseCommon.EmailAlert{
	DocID:        "",
	ParseFailure: abiCommon.ParseFailure{},
}

var (
	clients test.ClientMock
	ctx     *gin.Context
)

type TestSuite struct {
	suite.Suite
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (suite *TestSuite) SetupSuite() {
	expectedYear := 2022
	currentInstant := time.Date(expectedYear, 12, 01, 00, 00, 00, 0, time.UTC)
	// Make 'CurrentTime' return hard-coded time in tests
	CurrentTime = func() time.Time {
		return currentInstant
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = &http.Request{}

	c := gomock.NewController(suite.T())
	defer c.Finish()
	InitializeClientMocks(&clients, c)

}

func (suite *TestSuite) SetupTest() {

	requestData, err := GetTestRequest(&TestLog1)

	if err != nil {
		suite.Error(err, err.Error())
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestData))

}

func (suite *TestSuite) TestStoreFailureError() {

	ctrl := Controller{Clients: clients.Clients}

	newFailure := Failure{
		Timestamp: CurrentTime(),
		Log:       TestLog1,
		MessageID: "test",
		Tx:        TestLog1.TxHash.String(),
		LogIndex:  fmt.Sprint(TestLog1.Index),
		Contract:  TestLog1.Address.String(),
	}
	gomock.Eq(&newFailure)

	clients.Clients.EXPECT().GetMdb().Return(clients.Mdb.Client)
	clients.Mdb.Client.EXPECT().Database(gomock.Eq("")).Return(clients.Mdb.Database)
	clients.Mdb.Database.EXPECT().Collection(gomock.Eq("failures")).Return(clients.Mdb.Collection)
	clients.Mdb.Collection.EXPECT().InsertOne(gomock.Eq(ctx), gomock.Eq(&newFailure)).Return(nil, errors.New("failed to store failure"))

	ctrl.StoreFailure(ctx)
}

func (suite *TestSuite) TestStoreFailureSuccess() {

	ctrl := Controller{Clients: clients.Clients}

	newFailure := Failure{
		Timestamp: CurrentTime(),
		Log:       TestLog1,
		MessageID: "test",
		Tx:        TestLog1.TxHash.String(),
		LogIndex:  fmt.Sprint(TestLog1.Index),
		Contract:  TestLog1.Address.String(),
	}
	gomock.Eq(&newFailure)

	clients.Clients.EXPECT().GetMdb().Return(clients.Mdb.Client)
	clients.Mdb.Client.EXPECT().Database(gomock.Eq("")).Return(clients.Mdb.Database)
	clients.Mdb.Database.EXPECT().Collection(gomock.Eq("failures")).Return(clients.Mdb.Collection)
	clients.Mdb.Collection.EXPECT().InsertOne(gomock.Eq(ctx), gomock.Eq(&newFailure)).Return(&mongo.InsertOneResult{InsertedID: 0}, nil)

	ctrl.StoreFailure(ctx)
}
