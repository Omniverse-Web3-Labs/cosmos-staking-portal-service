package delegator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

type UserDelegationRequest struct {
	Address string `json:"address"`
}

func UserDelegationAction(c *gin.Context) {
	var data any
	var request UserDelegationRequest
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		delegation, err := service.DelegatorService.GetUserDelegations(c.Request.Context(), request.Address)

		if err != nil {
			return err
		}

		data = delegation
		return nil
	}()
	resp.Finish(c, err, data)
}
