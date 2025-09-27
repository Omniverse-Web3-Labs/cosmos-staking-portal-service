package validator

import (
	"app/internal/http/resp"
	"app/internal/model"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

func CommissionAction(c *gin.Context) {
	var data any
	var request model.ValidatorEventQuery
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		page, err := service.SubqueryService.GetValidatorEvents(c.Request.Context(), &request)
		if err != nil {
			return err
		}
		data = page.ToMap()
		return nil
	}()
	resp.Finish(c, err, data)
}
