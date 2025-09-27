package delegator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

func StakingParamsAction(c *gin.Context) {
	var data any
	err := func() error {
		params, err := service.DelegatorService.GetStakingParams(c.Request.Context())

		if err != nil {
			return err
		}

		data = params
		return nil
	}()
	resp.Finish(c, err, data)
}
