package delegator

import (
	"app/internal/http/resp"
	"app/internal/model"
	"app/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

func RewardRateAction(c *gin.Context) {
	var data any
	var request model.StakingRewardRateQuery
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		periodInt64, err := request.Period.Int64()
		if err != nil {
			return errors.New("invalid period")
		}
		periodInt := int(periodInt64)
		if periodInt != 1 && periodInt != 7 && periodInt != 30 {
			return errors.New("invalid period")
		}
		list, err := service.StakingRewardService.ListStakingRewardRate(c.Request.Context(), request.Delegator, periodInt)
		if err != nil {
			return err
		}

		data = list
		return nil
	}()
	resp.Finish(c, err, data)
}
