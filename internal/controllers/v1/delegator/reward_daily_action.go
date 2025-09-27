package delegator

import (
	"app/internal/http/resp"
	"app/internal/model"
	"app/internal/service"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

func RewardDailyAction(c *gin.Context) {
	var data any
	var request model.StakingRewardDailyQuery
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		rangeInt64, err := request.Range.Int64()
		if err != nil {
			return errors.New("invalid range")
		}
		rangeInt := int(rangeInt64)
		if rangeInt != 7 && rangeInt != 30 && rangeInt != 90 {
			return errors.New("invalid range")
		}
		now := time.Now()
		tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		dayTimes, beforeDayTimes, afterDayTimes, err := service.StakingRewardService.GetValidDayTimeByRange(c.Request.Context(), rangeInt, tm)
		if err != nil {
			return err
		}

		list, err := service.StakingRewardService.ListStakingRewardDaily(c.Request.Context(), request.Delegator, dayTimes)
		if err != nil {
			return err
		}
		emptyStakingRewardDailyList := service.StakingRewardService.CreateEmptyStakingRewardDailyList(c.Request.Context(), beforeDayTimes)
		stakingRewardDailyList := append(emptyStakingRewardDailyList, list...)
		emptyStakingRewardDailyList = service.StakingRewardService.CreateEmptyStakingRewardDailyList(c.Request.Context(), afterDayTimes)
		stakingRewardDailyList = append(stakingRewardDailyList, emptyStakingRewardDailyList...)
		data = stakingRewardDailyList
		return nil
	}()
	resp.Finish(c, err, data)
}
