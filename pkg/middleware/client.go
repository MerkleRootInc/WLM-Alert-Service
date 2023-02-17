package middleware

import (
	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/client"
	"github.com/MerkleRootInc/WLM-Event-Parser-Service/pkg/common"
	"github.com/gin-gonic/gin"
)

func InitClients(clients client.IClients, err common.ClientInitErr) gin.HandlerFunc {
	var (
		location = "Middleware.InitClients"
		errOpts  = common.EventParserErrOpts{
			Code:     &common.StatusBadRequest,
			Location: location,
		}
	)

	return func(c *gin.Context) {

		// return an internal server error if any clients had trouble initializing
		if err.Mdb != nil {
			common.RaiseAlertErr(c, errOpts, err.Mdb)
			c.Abort()

			return
		} else if err.Eth != nil {
			common.RaiseAlertErr(c, errOpts, err.Eth)
			c.Abort()

			return
		} else if err.Ps != nil {
			common.RaiseAlertErr(c, errOpts, err.Ps)
			c.Abort()

			return
		} else if err.Sm != nil {
			common.RaiseAlertErr(c, errOpts, err.Sm)
			c.Abort()

			return
		} else if err.Cs != nil {
			common.RaiseAlertErr(c, errOpts, err.Cs)
			c.Abort()

			return
		}

		c.Next()
	}
}
