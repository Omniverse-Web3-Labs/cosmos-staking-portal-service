package service

import (
	"app/internal/db/postgres"
	dbModel "app/internal/db/postgres/model"
	"context"
	"errors"
)

// 快照时间范围重放器
type RewardSnapshotTimetangeReplayer struct {
	validator               string
	delegator               string
	startTime               int64
	endTime                 int64
	dsIndex                 int
	delegatorPowerSnapshots *[]*dbModel.PowerSnapshots
	validatorRewardRecords  []*ValidatorReward
	resultRewardAmount      float64
}

func NewRewardSnapshotTimetangeReplayer(validator string, delegator string, startTime int64, endTime int64) *RewardSnapshotTimetangeReplayer {
	return &RewardSnapshotTimetangeReplayer{
		validator: validator,
		delegator: delegator,
		startTime: startTime,
		endTime:   endTime,
		dsIndex:   -1,
	}
}

func (s *RewardSnapshotTimetangeReplayer) GetResultRewardAmount() float64 {
	return s.resultRewardAmount
}

func (s *RewardSnapshotTimetangeReplayer) SetDelegatorPowerSnapshots(delegatorPowerSnapshots *[]*dbModel.PowerSnapshots) {
	s.delegatorPowerSnapshots = delegatorPowerSnapshots
}

func (s *RewardSnapshotTimetangeReplayer) SetValidatorRewardRecords(validatorRewardRecords []*ValidatorReward) {
	s.validatorRewardRecords = validatorRewardRecords
}

// 处理下一个快照索引
func (s *RewardSnapshotTimetangeReplayer) HandleNextIndex(timestamp int64) {
	for s.dsIndex < len(*s.delegatorPowerSnapshots)-1 && (*s.delegatorPowerSnapshots)[s.dsIndex+1].Timestamp < timestamp {
		s.dsIndex++
	}
}

// 处理下一个快照索引并记录矫正stakingAmount
func (s *RewardSnapshotTimetangeReplayer) HandleNextIndexAndCalCorrectAmount(timestamp int64) (correctStakingAmount float64) {
	for s.dsIndex < len(*s.delegatorPowerSnapshots)-1 && (*s.delegatorPowerSnapshots)[s.dsIndex+1].Timestamp < timestamp {
		s.dsIndex++
		timeDiff := (*s.delegatorPowerSnapshots)[s.dsIndex].Timestamp % DAY_TIME
		correctStakingAmount += (*s.delegatorPowerSnapshots)[s.dsIndex].Amount * float64(DAY_TIME-timeDiff) / float64(DAY_TIME)
	}
	return
}

func (s *RewardSnapshotTimetangeReplayer) GetCurrentDelegatorTotalShares() float64 {
	if s.dsIndex < 0 {
		return 0
	}
	return (*s.delegatorPowerSnapshots)[s.dsIndex].DelegatorTotalShares
}

func (s *RewardSnapshotTimetangeReplayer) Replay(ctx context.Context) error {
	resultRewardAmount := 0.0
	for _, validatorReward := range s.validatorRewardRecords {
		s.HandleNextIndexAndCalCorrectAmount(validatorReward.Timestamp)
		delegatorTotalShares := s.GetCurrentDelegatorTotalShares()
		amount := validatorReward.RewardAmount * delegatorTotalShares / validatorReward.TotalShare
		resultRewardAmount += amount
	}
	s.resultRewardAmount = resultRewardAmount
	return nil
}

func (s *RewardSnapshotTimetangeReplayer) Init(ctx context.Context) error {
	// 获取validator收益记录

	powerSnapshots, err := StakingRewardService.GetPowerSnapshotsByValidator(ctx, s.validator, s.startTime, s.endTime)
	if err != nil {
		return err
	}
	validatorRewardList, err := s.ListValidatorReward(ctx, s.validator, powerSnapshots, s.startTime, s.endTime)
	if err != nil {
		return err
	}
	s.validatorRewardRecords = validatorRewardList

	delegatorPowerSnapshots, err := s.GetPowerSnapshotsByDelegator(ctx, s.validator, s.delegator, s.startTime, s.endTime)
	if err != nil {
		return err
	}
	s.delegatorPowerSnapshots = &delegatorPowerSnapshots
	return nil
}

// 获取delegator的质押快照记录
func (s *RewardSnapshotTimetangeReplayer) GetPowerSnapshotsByDelegator(ctx context.Context, validator string, delegator string, startTime int64, endTime int64) ([]*dbModel.PowerSnapshots, error) {
	db := postgres.Subquery()
	list := make([]*dbModel.PowerSnapshots, 0)
	err := db.Model(&dbModel.PowerSnapshots{}).Where("delegator = ?", delegator).
		Where("timestamp >= ?", startTime).Where("timestamp < ?", endTime).Where("validator = ?", validator).
		Order("timestamp ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}

	var firstPowerSnapshots []*dbModel.PowerSnapshots = make([]*dbModel.PowerSnapshots, 0)
	err = db.Raw(`
		SELECT *
		FROM power_snapshots
		WHERE delegator = ?
		AND timestamp < ?
		AND validator = ?
		ORDER BY validator, timestamp DESC
		LIMIT 1
	`, delegator, startTime, validator).Find(&firstPowerSnapshots).Error
	if err != nil {
		return nil, errors.New("get power snapshots by delegator failed")
	}
	list = append(firstPowerSnapshots, list...)

	return list, nil
}

func (s *RewardSnapshotTimetangeReplayer) ListValidatorReward(ctx context.Context, validator string, powerSnapshots []*dbModel.PowerSnapshots, startTime int64, endTime int64) ([]*ValidatorReward, error) {
	validatorRewardList, err := s.ListValidatorRewardFromDB(ctx, validator, startTime, endTime)
	if err != nil {
		return nil, err
	}
	result := make([]*ValidatorReward, 0)
	psIndex := -1
	totalShare := 0.0
	var validatorRewardModel *ValidatorReward = nil
	for _, validatorReward := range validatorRewardList {
		for psIndex < len(powerSnapshots)-1 && validatorReward.Timestamp > powerSnapshots[psIndex+1].Timestamp {
			psIndex++
			totalShare = powerSnapshots[psIndex].ValidatorTotalShares
			validatorRewardModel = nil
		}
		if validatorRewardModel == nil {
			validatorRewardModel = &ValidatorReward{
				RewardAmount: validatorReward.RewardAmount - validatorReward.CommissionAmount,
				Timestamp:    validatorReward.Timestamp,
				TotalShare:   totalShare,
			}
			result = append(result, validatorRewardModel)
		} else {
			validatorRewardModel.RewardAmount += validatorReward.RewardAmount - validatorReward.CommissionAmount
		}
	}
	return result, nil
}

func (s *RewardSnapshotTimetangeReplayer) ListValidatorRewardFromDB(ctx context.Context, validator string, startTime int64, endTime int64) ([]*dbModel.ValidatorRewardRecords, error) {
	list := make([]*dbModel.ValidatorRewardRecords, 0)

	db := postgres.Subquery()

	err := db.Model(&dbModel.ValidatorRewardRecords{}).
		Where("validator = ?", validator).Where("timestamp >= ?", startTime).
		Where("timestamp < ?", endTime).Order("timestamp ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
