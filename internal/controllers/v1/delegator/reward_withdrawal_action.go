package delegator

import (
	"app/internal/http/resp"
	"app/internal/model"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

func RewardWithdrawalAction(c *gin.Context) {
	var data any
	var request model.StakingRewardWithdrawalQuery
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		list, err := service.StakingRewardService.ListStakingRewardWithdrawal(c.Request.Context(), &request)
		if err != nil {
			return err
		}

		data = list
		return nil
	}()
	resp.Finish(c, err, data)
}
