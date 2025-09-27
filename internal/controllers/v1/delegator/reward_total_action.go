package delegator

import (
	"app/internal/http/resp"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

type StakingRewardTotalQuery struct {
	Delegator string `json:"delegator"`
}

func RewardTotalAction(c *gin.Context) {
	var data any
	var request StakingRewardTotalQuery
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		totalReward, err := service.StakingRewardService.GetDelegatorTotalReward(c.Request.Context(), request.Delegator)
		if err != nil {
			return err
		}
		data = gin.H{
			"totalReward": totalReward,
		}
		return nil
	}()
	resp.Finish(c, err, data)
}
