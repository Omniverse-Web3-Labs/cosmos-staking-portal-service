package validator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

type ValidatorBatchRequest struct {
	Addresses []string `json:"addresses"`
}

func BatchAction(c *gin.Context) {
	var data any
	var request ValidatorBatchRequest
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		validators, err := service.ValidatorService.GetValidators(c.Request.Context(), request.Addresses)
		if err != nil {
			return err
		}
		data = validators

		return nil
	}()
	resp.Finish(c, err, data)
}
