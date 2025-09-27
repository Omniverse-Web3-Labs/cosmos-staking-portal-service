package delegator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

type DelegatorDetailRequest struct {
	Address string `json:"address"`
}

func DetailAction(c *gin.Context) {
	var data any
	var request DelegatorDetailRequest
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		delegator, err := service.DelegatorService.GetDelegatorDetail(c.Request.Context(), request.Address)

		if err != nil {
			return err
		}

		data = delegator
		return nil
	}()
	resp.Finish(c, err, data)
}
