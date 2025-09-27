package validator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

type ValidatorDetailRequest struct {
	Address string `json:"address"`
}

type ValidatorDetailResponse struct {
	Detail service.ValidatorDetail `json:"detail"`
}

func DetailAction(c *gin.Context) {
	var data any
	var request ValidatorDetailRequest
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		validator, err := service.ValidatorService.GetValidatorDetail(c.Request.Context(), request.Address)

		if err != nil {
			return err
		}

		data = validator
		return nil
	}()
	resp.Finish(c, err, data)
}
