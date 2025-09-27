package validator

import (
	"app/internal/http/resp"
	"app/internal/model"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

func PowerEventAction(c *gin.Context) {
	var data any
	var request model.PowerEventQuery
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		page, err := service.SubqueryService.GetPowerEvents(c.Request.Context(), &request)
		if err != nil {
			return err
		}
		data = page
		return nil
	}()
	resp.Finish(c, err, data)
}
