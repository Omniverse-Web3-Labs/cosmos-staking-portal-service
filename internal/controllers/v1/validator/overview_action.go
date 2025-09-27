package validator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

func OverviewAction(c *gin.Context) {
	var data any
	err := func() error {
		overview, err := service.ValidatorService.GetStakeOverview(c.Request.Context())
		if err != nil {
			return err
		}

		data = overview

		return nil
	}()
	resp.Finish(c, err, data)
}
