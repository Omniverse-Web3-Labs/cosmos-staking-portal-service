package service

import (
	dbModel "app/internal/db/postgres/model"
	"app/internal/model"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGetDayTimeByRange(t *testing.T) {
	tm := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	dayTime := StakingRewardService.GetDayTimeByRange(7, tm)

	if dayTime[0] != 1682294400000 {
		t.Fatal("dayTime[0] is not 1682294400000")
	}
	if dayTime[6] != 1682812800000 {
		t.Fatal("dayTime[6] is not 1682812800000")
	}
}

func TestGetStakingRewardDaily(t *testing.T) {
	tm := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
	dayTimes := StakingRewardService.GetDayTimeByRange(7, tm)

	ctx := context.Background()

	stakingRewardDaily, err := StakingRewardService.GetStakingRewardDailyFromDB(ctx, "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg", dayTimes)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	t.Logf("stakingRewardDaily:%+v\n", stakingRewardDaily)
}

func TestGetPowerSnapshotsByDelegator(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayTimes := StakingRewardService.GetDayTimeByRange(7, tm)

	powerSnapshots, err := StakingRewardService.GetPowerSnapshotsByDelegator(ctx, "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg", dayTimes)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	t.Logf("powerSnapshots:%+v\n", powerSnapshots)
}

// GetValidatorByDelegator
func TestGetValidatorByDelegator(t *testing.T) {
	ctx := context.Background()
	validators, err := StakingRewardService.GetValidatorByDelegator(ctx, "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg")
	if err != nil {
		t.Fatal("error", err.Error())
	}
	t.Logf("validators:%+v\n", validators)
}

// GetValidatorRewardRecordMap
func TestGetValidatorRewardRecordMap(t *testing.T) {
	ctx := context.Background()
	validator := "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg"
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayTimes := StakingRewardService.GetDayTimeByRange(7, tm)

	validatorRewardRecordMap, err := StakingRewardService.GetValidatorRewardRecordMap(ctx, validator, dayTimes)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	for k, _ := range validatorRewardRecordMap {
		fmt.Println("key:", k)

	}
}

// Replayer
func TestReplayer(t *testing.T) {
	ctx := context.Background()
	delegator := "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg"
	validator := "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg"
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayTimes := StakingRewardService.GetDayTimeByRange(7, tm)

	replayer := NewSnapshotReplayer(dayTimes, delegator, validator)
	delegatorPowerSnapshots := []*dbModel.PowerSnapshots{
		{
			Delegator:            delegator,
			Validator:            validator,
			Timestamp:            dayTimes[0] + 1,
			Amount:               100,
			DelegatorTotalShares: 100,
			ValidatorTotalShares: 100,
			ValidatorTotalAmount: 100,
		},
		{
			Delegator:            delegator,
			Validator:            validator,
			Timestamp:            dayTimes[0] + DAY_TIME/2,
			Amount:               100,
			DelegatorTotalShares: 200,
			ValidatorTotalShares: 200,
			ValidatorTotalAmount: 200,
		},
	}
	validatorRewardRecords := make(map[int64][]*ValidatorReward)
	validatorRewardRecords[dayTimes[0]] = []*ValidatorReward{
		{
			RewardAmount: 100,
			Timestamp:    dayTimes[0] + 1,
			TotalShare:   100,
		},
		{
			RewardAmount: 100,
			Timestamp:    dayTimes[0] + DAY_TIME/2,
			TotalShare:   200,
		},
	}
	validatorRewardRecords[dayTimes[1]] = []*ValidatorReward{
		{
			RewardAmount: 200,
			Timestamp:    dayTimes[1] + 1,
			TotalShare:   200,
		},
	}
	validatorRewardRecords[dayTimes[6]] = []*ValidatorReward{
		{
			RewardAmount: 200,
			Timestamp:    dayTimes[6] + 120,
			TotalShare:   200,
		},
	}
	replayer.SetDelegatorPowerSnapshots(&delegatorPowerSnapshots)
	replayer.SetValidatorRewardRecords(&validatorRewardRecords)

	stakingRewardDailyList, err := replayer.Replay(ctx, delegator)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	first := stakingRewardDailyList[0]
	second := stakingRewardDailyList[1]
	last := stakingRewardDailyList[len(stakingRewardDailyList)-1]
	if first.StakingAmount >= 100 && first.StakingAmount <= 99.99 {
		t.Fatal("first.StakingAmount is not 100")
	}
	if first.RewardAmount != 50 {
		t.Fatal("first.RewardAmount is not 50")
	}
	if second.StakingAmount != 200 {
		t.Fatal("second.StakingAmount is not 200")
	}
	if second.RewardAmount != 200 {
		t.Fatal("second.RewardAmount is not 200")
	}
	if last.StakingAmount != 200 {
		t.Fatal("last.StakingAmount is not 200")
	}
	if last.RewardAmount != 200 {
		t.Fatal("last.RewardAmount is not 200")
	}

	t.Logf("stakingRewardDailyList:%+v\n", stakingRewardDailyList)
}

func TestCalculateAndStoreStakingReward(t *testing.T) {
	ctx := context.Background()
	delegator := "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg"
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayTimes := StakingRewardService.GetDayTimeByRange(7, tm)

	_, err := StakingRewardService.CalculateAndStoreStakingReward(ctx, delegator, dayTimes)
	if err != nil {
		t.Fatal("error", err.Error())
	}
}

// ListStakingRewardRate
func TestListStakingRewardRate(t *testing.T) {
	ctx := context.Background()
	delegator := "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg"
	rates, err := StakingRewardService.ListStakingRewardRate(ctx, delegator, 1)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	t.Logf("rates:%+v\n", rates)
}

func TestListStakingRewardWithdrawal(t *testing.T) {
	ctx := context.Background()
	query := &model.StakingRewardWithdrawalQuery{
		Delegator: "omni195tf7tuvxagrsq7seh8579qrjvxc6j9wgnpz86",
	}
	paginator, err := StakingRewardService.ListStakingRewardWithdrawal(ctx, query)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	t.Logf("paginator:%+v\n", paginator)
}

func TestGetDelegatorTotalReward(t *testing.T) {
	ctx := context.Background()
	delegator := "omni195tf7tuvxagrsq7seh8579qrjvxc6j9wgnpz86"
	totalReward, err := StakingRewardService.GetDelegatorTotalReward(ctx, delegator)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	t.Logf("totalReward:%+v\n", totalReward)
}

func TestCompressValidatorRewardByDay(t *testing.T) {
	ctx := context.Background()
	dayTime := int64(1731024000000)
	err := StakingRewardService.CompressValidatorRewardByDay(ctx, dayTime)
	if err != nil {
		t.Fatal("error", err.Error())
	}
}
