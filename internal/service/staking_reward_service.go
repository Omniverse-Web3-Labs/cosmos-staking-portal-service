package service

import (
	"app/internal/cache/redisclient"
	"app/internal/db/postgres"
	dbModel "app/internal/db/postgres/model"
	"app/internal/model"
	"app/pkg/utils"
	"encoding/json"
	"fmt"

	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//	质押收益设计 ListStakingRewardDaily

// 1.通过dayTimes从db查询delegator每日质押收益记录(StakingRewardDaily)并计算空缺数据startTime和endTime
// 2 查询delegator质押过的validator,如果validator为空则返回空数据
// 3.通过startTime和endTime查询delegator的质押快照记录(PowerSnapshots)
// 4.通过startTime和endTime查询validator的质押快照记录(PowerSnapshots)
// 5.通过startTime和endTime查询validator收益记录(ValidatorReward)
// 6.通过validator的质押快照记录和validator的收益记录更新收益记录的totalShare
// 7.通过validator的收益记录和delegator的质押快照记录构建快照回放器(SnapshotReplayer)
// 8.回放计算delegator的质押每日收益

// 一天毫秒数
const DAY_TIME = 1000 * 60 * 60 * 24
const CHAIN_START_TIME = 1725580800000

// redis缓存key, validator 收益记录已经压缩的最新时间戳
const VALIDATOR_REWARD_COMPRESSED_LATEST_TIME_CACHE_KEY = "validator_reward:compressed:latest_time"
const VALIDATOR_REWARD_RECORDS_CACHE_KEY = "validator_reward:records:%s:%d"

// 快照回放器
type SnapshotReplayer struct {
	validator string
	delegator string
	dayTimes  []int64
	//delegatorSnapshotIndex
	dsIndex                 int
	delegatorPowerSnapshots *[]*dbModel.PowerSnapshots
	validatorRewardRecords  *map[int64][]*ValidatorReward
}

func NewSnapshotReplayer(dayTimes []int64, delegator string, validator string) *SnapshotReplayer {
	return &SnapshotReplayer{
		dayTimes:  dayTimes,
		delegator: delegator,
		validator: validator,
		dsIndex:   -1,
	}
}

func (s *SnapshotReplayer) SetDelegatorPowerSnapshots(delegatorPowerSnapshots *[]*dbModel.PowerSnapshots) {
	s.delegatorPowerSnapshots = delegatorPowerSnapshots
}

func (s *SnapshotReplayer) SetValidatorRewardRecords(validatorRewardRecords *map[int64][]*ValidatorReward) {
	s.validatorRewardRecords = validatorRewardRecords
}

// 处理下一个快照索引
func (s *SnapshotReplayer) HandleNextIndex(timestamp int64) {
	for s.dsIndex < len(*s.delegatorPowerSnapshots)-1 && (*s.delegatorPowerSnapshots)[s.dsIndex+1].Timestamp < timestamp {
		s.dsIndex++
	}
}

// 处理下一个快照索引并记录矫正stakingAmount
func (s *SnapshotReplayer) HandleNextIndexAndCalCorrectAmount(timestamp int64) (correctStakingAmount float64) {
	for s.dsIndex < len(*s.delegatorPowerSnapshots)-1 && (*s.delegatorPowerSnapshots)[s.dsIndex+1].Timestamp < timestamp {
		s.dsIndex++
		timeDiff := (*s.delegatorPowerSnapshots)[s.dsIndex].Timestamp % DAY_TIME
		correctStakingAmount += (*s.delegatorPowerSnapshots)[s.dsIndex].Amount * float64(DAY_TIME-timeDiff) / float64(DAY_TIME)
	}
	return
}

func (s *SnapshotReplayer) GetCurrentStakingAmount() float64 {
	if s.dsIndex < 0 {
		return 0
	}
	powerSnapshot := (*s.delegatorPowerSnapshots)[s.dsIndex]
	return powerSnapshot.ValidatorTotalAmount * powerSnapshot.DelegatorTotalShares / powerSnapshot.ValidatorTotalShares
}

func (s *SnapshotReplayer) GetCurrentDelegatorTotalShares() float64 {
	if s.dsIndex < 0 {
		return 0
	}
	return (*s.delegatorPowerSnapshots)[s.dsIndex].DelegatorTotalShares
}

func (s *SnapshotReplayer) Replay(ctx context.Context, delegator string) ([]*dbModel.StakingRewardDaily, error) {
	stakingRewardDailyList := make([]*dbModel.StakingRewardDaily, 0)
	for _, dayTime := range s.dayTimes {
		s.HandleNextIndex(dayTime)
		stakingRewardDaily := &dbModel.StakingRewardDaily{
			Delegator:     delegator,
			DayTime:       dayTime,
			RewardAmount:  0,
			StakingAmount: s.GetCurrentStakingAmount(),
			CreatedAt:     time.Now().UnixMilli(),
		}
		stakingRewardDailyList = append(stakingRewardDailyList, stakingRewardDaily)
		if s.delegatorPowerSnapshots == nil || len(*s.delegatorPowerSnapshots) == 0 {
			continue
		}
		validatorRewardList := (*s.validatorRewardRecords)[dayTime]
		if validatorRewardList == nil {
			continue
		}
		//累计矫正的stakingAmount值
		correctStakingAmount := 0.0
		for _, validatorReward := range validatorRewardList {
			correctStakingAmount += s.HandleNextIndexAndCalCorrectAmount(validatorReward.Timestamp)
			stakingRewardDaily.RewardAmount += validatorReward.RewardAmount * s.GetCurrentDelegatorTotalShares() / validatorReward.TotalShare
		}
		stakingRewardDaily.StakingAmount += correctStakingAmount
	}
	return stakingRewardDailyList, nil
}

type stakingRewardService struct {
}

type ValidatorReward struct {
	RewardAmount float64
	Timestamp    int64
	TotalShare   float64
}

var StakingRewardService stakingRewardService

func (s *stakingRewardService) ListStakingRewardDaily(ctx context.Context, delegator string, dayTimes []int64) ([]*model.StakingRewardDaily, error) {
	rewardDailyList, err := s.GetStakingRewardDailyFromDB(ctx, delegator, dayTimes)
	if err != nil {
		return nil, err
	}
	//预期length
	if len(dayTimes) != len(rewardDailyList) {
		// 将rewardDailyList 转换为map
		rewardDailyMap := make(map[int64]bool)
		for _, rewardDaily := range rewardDailyList {
			rewardDailyMap[rewardDaily.DayTime] = true
		}
		// 计算需要计算的dayTimes
		needCalculateDayTimes := make([]int64, 0)
		for _, dayTime := range dayTimes {
			if _, ok := rewardDailyMap[dayTime]; !ok {
				needCalculateDayTimes = append(needCalculateDayTimes, dayTime)
			}
		}
		if len(needCalculateDayTimes) > 0 {
			// 计算全部days收益
			if _, err = s.CalculateAndStoreStakingReward(ctx, delegator, needCalculateDayTimes); err != nil {
				return nil, err
			}
		}
		// 获取计算后的收益
		rewardDailyList, err = s.GetStakingRewardDailyFromDB(ctx, delegator, dayTimes)
		if err != nil {
			return nil, err
		}
	}

	list := make([]*model.StakingRewardDaily, len(dayTimes))
	// 填充数据
	rewardDailyMap := make(map[int64]*dbModel.StakingRewardDaily)
	for _, rewardDaily := range rewardDailyList {
		rewardDailyMap[rewardDaily.DayTime] = rewardDaily
	}
	for index, dayTime := range dayTimes {
		rewardDaily, ok := rewardDailyMap[dayTime]
		if ok {
			list[index] = &model.StakingRewardDaily{
				Timestamp:     rewardDaily.DayTime,
				RewardAmount:  rewardDaily.RewardAmount,
				StakingAmount: rewardDaily.StakingAmount,
			}
		} else {
			list[index] = &model.StakingRewardDaily{
				Timestamp:     dayTime,
				RewardAmount:  0,
				StakingAmount: 0,
			}
		}
	}

	return list, nil
}

// 计算收益并存储到数据库
func (s *stakingRewardService) CalculateAndStoreStakingReward(ctx context.Context, delegator string, dayTimes []int64) ([]*dbModel.StakingRewardDaily, error) {
	//查询所有质押的validator
	validators, err := s.GetValidatorByDelegator(ctx, delegator)
	if err != nil {
		return nil, err
	}
	if len(validators) == 0 {
		return nil, nil
	}
	// 获取delegator 质押快照记录
	powerSnapshotsMap, err := s.GetPowerSnapshotsByDelegator(ctx, delegator, dayTimes)
	if err != nil {
		return nil, err
	}

	allStakingRewardDailyList := make([]*dbModel.StakingRewardDaily, len(dayTimes))

	for _, validator := range validators {
		//获取validator的收益记录
		validatorRewardMap, err := s.GetValidatorRewardRecordMap(ctx, validator, dayTimes)
		if err != nil {
			return nil, err
		}
		//获取delegator的质押记录
		delegatorPowerSnapshots := powerSnapshotsMap[validator]
		//创建快照回放器
		snapshotReplayer := NewSnapshotReplayer(dayTimes, delegator, validator)
		snapshotReplayer.SetDelegatorPowerSnapshots(&delegatorPowerSnapshots)
		snapshotReplayer.SetValidatorRewardRecords(&validatorRewardMap)
		//回放
		stakingRewardDailyList, err := snapshotReplayer.Replay(ctx, delegator)
		if err != nil {
			return nil, err
		}
		for index, stakingRewardDaily := range stakingRewardDailyList {
			if allStakingRewardDailyList[index] == nil {
				allStakingRewardDailyList[index] = stakingRewardDaily
			} else {
				allStakingRewardDailyList[index].RewardAmount += stakingRewardDaily.RewardAmount
				allStakingRewardDailyList[index].StakingAmount += stakingRewardDaily.StakingAmount
			}
		}
	}
	// 批量存储到数据库,忽略重复的记录
	db := postgres.Default()
	err = db.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(allStakingRewardDailyList, 100).Error
	if err != nil {
		return nil, err
	}
	return allStakingRewardDailyList, nil
}

// 从数据库获取每日质押收益记录
func (s *stakingRewardService) GetStakingRewardDailyFromDB(ctx context.Context, delegator string, dayTime []int64) ([]*dbModel.StakingRewardDaily, error) {
	list := make([]*dbModel.StakingRewardDaily, 0)
	db := postgres.Default()
	err := db.Model(&dbModel.StakingRewardDaily{}).WithContext(ctx).Where("delegator = ?", delegator).Where("day_time IN (?)", dayTime).Order("day_time ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 获取validator收益记录
func (s *stakingRewardService) GetValidatorRewardRecordMap(ctx context.Context, validator string, dayTimes []int64) (map[int64][]*ValidatorReward, error) {

	validatorRewardMap := make(map[int64][]*ValidatorReward)
	for _, dayTime := range dayTimes {
		dayTimeList, err := s.GetValidatorRewardFromCache(ctx, validator, dayTime)
		if err != nil {
			return nil, err
		}
		validatorRewardMap[dayTime] = dayTimeList
	}

	return validatorRewardMap, nil
}

func (s *stakingRewardService) GetValidatorRewardFromCache(ctx context.Context, validator string, dayTime int64) ([]*ValidatorReward, error) {
	cacheKey := fmt.Sprintf(VALIDATOR_REWARD_RECORDS_CACHE_KEY, validator, dayTime)
	cacheValue := redisclient.RedisClient.Get(ctx, cacheKey)
	err := cacheValue.Err()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, err
	}
	cacheValueBytes, _ := cacheValue.Bytes()
	list := make([]*ValidatorReward, 0)
	err = json.Unmarshal(cacheValueBytes, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *stakingRewardService) SetValidatorRewardToCache(ctx context.Context, validator string, dayTime int64, list []*ValidatorReward) error {
	cacheKey := fmt.Sprintf(VALIDATOR_REWARD_RECORDS_CACHE_KEY, validator, dayTime)
	cacheValue, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = redisclient.RedisClient.Set(ctx, cacheKey, cacheValue, time.Hour*24*365).Err()
	if err != nil {
		return err
	}
	return nil
}
func (s *stakingRewardService) ListValidatorReward(ctx context.Context, validator string, dayTime int64) ([]*dbModel.ValidatorRewardRecords, error) {
	list := make([]*dbModel.ValidatorRewardRecords, 0)

	db := postgres.Subquery()

	endTime := dayTime + DAY_TIME
	err := db.Model(&dbModel.ValidatorRewardRecords{}).WithContext(ctx).
		Where("validator = ?", validator).Where("timestamp >= ?", dayTime).
		Where("timestamp < ?", endTime).Order("timestamp ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 获取validator的质押快照记录
func (s *stakingRewardService) GetPowerSnapshotsByValidator(ctx context.Context, validator string, startTime int64, endTime int64) ([]*dbModel.PowerSnapshots, error) {
	db := postgres.Subquery()
	firstPowerSnapshot := make([]*dbModel.PowerSnapshots, 0)

	//查询时间外validator的第一个快照
	err := db.Model(&dbModel.PowerSnapshots{}).WithContext(ctx).Where("validator = ?", validator).
		Where("timestamp < ?", startTime).
		Order("timestamp DESC").Limit(1).Find(&firstPowerSnapshot).Error
	if err != nil {
		return nil, err
	}

	//查询时间内的所有快照
	list := make([]*dbModel.PowerSnapshots, 0)
	db.Model(&dbModel.PowerSnapshots{}).WithContext(ctx).Where("validator = ?", validator).
		Where("timestamp >= ?", startTime).Where("timestamp < ?", endTime).
		Order("timestamp ASC").Find(&list)
	//合并快照
	list = append(firstPowerSnapshot, list...)

	return list, nil
}

// 获取delegator的质押快照记录
func (s *stakingRewardService) GetPowerSnapshotsByDelegator(ctx context.Context, delegator string, dayTimes []int64) (map[string][]*dbModel.PowerSnapshots, error) {
	db := postgres.Subquery()
	startTime := dayTimes[0]
	endTime := dayTimes[len(dayTimes)-1] + DAY_TIME
	list := make([]*dbModel.PowerSnapshots, 0)
	err := db.Model(&dbModel.PowerSnapshots{}).WithContext(ctx).Where("delegator = ?", delegator).
		Where("timestamp >= ?", startTime).Where("timestamp < ?", endTime).
		Order("timestamp ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}

	var firstPowerSnapshots []*dbModel.PowerSnapshots = make([]*dbModel.PowerSnapshots, 0)
	err = db.Raw(`
		SELECT DISTINCT ON (validator) *
		FROM power_snapshots
		WHERE delegator = ?
		AND timestamp < ?
		ORDER BY validator, timestamp DESC
	`, delegator, startTime).Find(&firstPowerSnapshots).Error
	if err != nil {
		return nil, errors.New("get power snapshots by delegator failed")
	}
	powerSnapshotsMap := make(map[string][]*dbModel.PowerSnapshots)
	list = append(firstPowerSnapshots, list...)
	for _, powerSnapshot := range list {
		powerSnapshotsMap[powerSnapshot.Validator] = append(powerSnapshotsMap[powerSnapshot.Validator], powerSnapshot)
	}
	return powerSnapshotsMap, nil
}

// 通过delegator获取validators
func (service *stakingRewardService) GetValidatorByDelegator(ctx context.Context, delegator string) ([]string, error) {
	validators := make([]string, 0)
	db := postgres.Subquery()
	err := db.Model(&dbModel.PowerSnapshots{}).WithContext(ctx).Where("delegator = ?", delegator).Group("validator").Pluck("validator", &validators).Error
	if err != nil {
		return nil, err
	}
	return validators, nil
}

// 通过range 获取时间数组
func (s *stakingRewardService) GetDayTimeByRange(timeRange int, time time.Time) []int64 {
	dayTime := make([]int64, timeRange)
	startTime := time.AddDate(0, 0, -timeRange)
	for i := 0; i < timeRange; i++ {
		dayTime[i] = startTime.AddDate(0, 0, i).UnixMilli()
	}
	return dayTime
}

// 获取有效时间数组
func (s *stakingRewardService) GetValidDayTimeByRange(ctx context.Context, timeRange int, time time.Time) ([]int64, []int64, []int64, error) {
	dayTimes := make([]int64, 0)
	beforeDayTimes := make([]int64, 0)
	afterDayTimes := make([]int64, 0)
	maxBlockTime, err := s.GetValidatorRewardCompressedLatestTime(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	startTime := time.AddDate(0, 0, -timeRange)
	for i := 0; i < timeRange; i++ {
		dayTime := startTime.AddDate(0, 0, i).UnixMilli()
		if dayTime >= CHAIN_START_TIME && dayTime <= maxBlockTime {
			dayTimes = append(dayTimes, dayTime)
		} else if dayTime < CHAIN_START_TIME { //链启动前的数据
			beforeDayTimes = append(beforeDayTimes, dayTime)
		} else { //未同步到的数据
			afterDayTimes = append(afterDayTimes, dayTime)
		}
	}
	return dayTimes, beforeDayTimes, afterDayTimes, nil
}

// 查询validator收益记录已经压缩的最新时间戳
func (s *stakingRewardService) GetValidatorRewardCompressedLatestTime(ctx context.Context) (int64, error) {

	cacheValue := redisclient.RedisClient.Get(ctx, VALIDATOR_REWARD_COMPRESSED_LATEST_TIME_CACHE_KEY)
	if cacheValue.Err() != nil && cacheValue.Err().Error() != "redis: nil" {
		return 0, cacheValue.Err()
	}
	var blockTime int64
	cacheValueBytes, _ := cacheValue.Bytes()
	if len(cacheValueBytes) > 0 {
		err := json.Unmarshal(cacheValueBytes, &blockTime)
		if err != nil {
			return 0, err
		}
		return blockTime, nil
	}

	return blockTime, nil
}

// 获取质押收益率
func (s *stakingRewardService) ListStakingRewardRate(ctx context.Context, delegator string, period int) ([]*model.StakingRewardRate, error) {
	var dayTimes []int64
	var beforeChainStartTime []int64
	var afterChainStartTime []int64
	var err error
	switch period {
	case 1:
		tm := utils.TimeUtilUTC.TodayStartTime()
		dayTimes, beforeChainStartTime, afterChainStartTime, err = s.GetValidDayTimeByRange(ctx, 7, tm)
		if err != nil {
			return nil, err
		}
	case 7:
		tm := utils.TimeUtilUTC.WeekStartTime()
		dayTimes, beforeChainStartTime, afterChainStartTime, err = s.GetValidDayTimeByRange(ctx, 7*7, tm)
		if err != nil {
			return nil, err
		}
	case 30:
		tm := utils.TimeUtilUTC.MonthStartTime()
		dayTimes, beforeChainStartTime, afterChainStartTime, err = s.GetValidDayTimeByRange(ctx, 30*7, tm)
		if err != nil {
			return nil, err
		}
	}
	stakingRewardDailyList, err := s.ListStakingRewardDaily(ctx, delegator, dayTimes)
	if err != nil {
		return nil, err
	}
	if len(stakingRewardDailyList) == 0 {
		return make([]*model.StakingRewardRate, 0), nil
	}
	emptyStakingRewardDailyList := s.CreateEmptyStakingRewardDailyList(ctx, beforeChainStartTime)
	stakingRewardDailyList = append(emptyStakingRewardDailyList, stakingRewardDailyList...)
	emptyStakingRewardDailyList = s.CreateEmptyStakingRewardDailyList(ctx, afterChainStartTime)
	stakingRewardDailyList = append(stakingRewardDailyList, emptyStakingRewardDailyList...)
	list := make([]*model.StakingRewardRate, 7)
	accApy := 0.0
	accApyCount := 0
	j := 0
	for index, stakingRewardDaily := range stakingRewardDailyList {
		//index余数
		mod := index % period
		if mod == 0 {
			accApy = 0
			accApyCount = 0
			list[j] = &model.StakingRewardRate{
				Timestamp: stakingRewardDaily.Timestamp,
			}
			j++
		}
		if stakingRewardDaily.StakingAmount != 0 {
			accApy += (stakingRewardDaily.RewardAmount / stakingRewardDaily.StakingAmount)
			accApyCount++
		}
		if mod == period-1 {
			if accApyCount != 0 {
				list[j-1].Rate = accApy * 365 * 10000 / float64(accApyCount)
			}
		}
	}
	return list, nil
}

func (s *stakingRewardService) CreateEmptyStakingRewardDailyList(ctx context.Context, dayTimes []int64) []*model.StakingRewardDaily {
	emptyStakingRewardDailyList := make([]*model.StakingRewardDaily, len(dayTimes))
	for i, dayTime := range dayTimes {
		emptyStakingRewardDailyList[i] = &model.StakingRewardDaily{
			Timestamp: dayTime,
		}
	}
	return emptyStakingRewardDailyList
}

// 获取质押提现记录
func (s *stakingRewardService) ListStakingRewardWithdrawal(ctx context.Context, query *model.StakingRewardWithdrawalQuery) (*utils.Paginator, error) {
	db := postgres.Subquery()
	if query.Delegator == "" {
		paginator := utils.NewPaginator(query.Current, query.Size)
		return paginator, nil
	}
	paginator := utils.NewPaginator(query.Current, query.Size)
	list := make([]*model.RewardWithdrawalEvents, 0)
	err := db.Model(&dbModel.RewardWithdrawalEvents{}).WithContext(ctx).
		Where("delegator = ?", query.Delegator).
		Select("reward_withdrawal_events.*").
		Order("timestamp DESC").Limit(paginator.Size + 1).Offset(paginator.Offset()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	paginator.HasMore = len(list) > paginator.Size
	if paginator.HasMore {
		list = list[:len(list)-1]
	}
	for _, rewardWithdrawalEvent := range list {
		validatorDescription, err := ValidatorDescriptionService.GetValidatorDescription(ctx, rewardWithdrawalEvent.Validator)
		if err != nil {
			return nil, err
		}
		rewardWithdrawalEvent.Picture = validatorDescription.Avatar
		rewardWithdrawalEvent.Moniker = validatorDescription.Moniker
	}
	paginator.SetList(list)
	return paginator, nil
}

// Delegator累计收益
func (s *stakingRewardService) GetDelegatorTotalReward(ctx context.Context, delegator string) (float64, error) {
	if delegator == "" {
		return 0, nil
	}
	db := postgres.Subquery()
	var totalReward float64
	err := db.Model(&dbModel.RewardWithdrawalEvents{}).WithContext(ctx).Select("COALESCE(SUM(amount), 0) as total_reward").
		Where("delegator = ?", delegator).Scan(&totalReward).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	return totalReward, nil
}

// 查询ValidatorRewardRecords表中最大的timestamp
func (s *stakingRewardService) GetValidatorRewardRecordsMaxTimestamp(ctx context.Context) (int64, error) {
	db := postgres.Subquery()
	var maxTimestamp int64
	err := db.Model(&dbModel.ValidatorRewardRecords{}).WithContext(ctx).Select("MAX(timestamp) as max_timestamp").Scan(&maxTimestamp).Error
	if err != nil {
		return 0, err
	}
	return maxTimestamp, nil
}

// 启动压缩validator收益记录任务
func (s *stakingRewardService) StartValidatorRewardCompressedTask(ctx context.Context) error {
	//查询validator收益记录已经压缩的最新时间戳
	latestTime, err := s.GetValidatorRewardCompressedLatestTime(ctx)
	if err != nil && err.Error() != "redis: nil" {
		return err
	}
	//若latestTime为0或者为空，从CHAIN_START_TIME开始压缩
	if latestTime == 0 {
		latestTime = CHAIN_START_TIME
	}
	// maxTimestamp 是昨天0点的时间戳
	maxTimestamp := utils.TimeUtilUTC.TodayStartTime().UnixMilli() - DAY_TIME
	dbMaxTimestamp, err := s.GetValidatorRewardRecordsMaxTimestamp(ctx)
	if err != nil {
		return err
	}
	if dbMaxTimestamp < maxTimestamp {
		maxTimestamp = dbMaxTimestamp / DAY_TIME * DAY_TIME
	}
	//若maxTimestamp小于latestTime，则从latestTime开始压缩
	if maxTimestamp <= latestTime {
		return nil
	}
	// 按天处理
	for dayTime := latestTime + DAY_TIME; dayTime <= maxTimestamp; dayTime += DAY_TIME {
		err = s.CompressValidatorRewardByDay(ctx, dayTime)
		if err != nil {
			return err
		}
		// 更新缓存
		redisclient.RedisClient.Set(ctx, VALIDATOR_REWARD_COMPRESSED_LATEST_TIME_CACHE_KEY, dayTime, -1)
	}
	return nil
}

// 按天压缩validator收益记录
func (s *stakingRewardService) CompressValidatorRewardByDay(ctx context.Context, dayTime int64) error {
	db := postgres.Subquery()
	// group by validator, timestamp
	validatorList := make([]string, 0)
	err := db.Model(&dbModel.ValidatorRewardRecords{}).WithContext(ctx).
		Where("timestamp >= ?", dayTime).Where("timestamp < ?", dayTime+DAY_TIME).
		Group("validator").Select("validator").Find(&validatorList).Error
	if err != nil {
		return err
	}
	for _, validator := range validatorList {
		err = s.CompressValidatorReward(ctx, validator, dayTime)
		if err != nil {
			return err
		}
	}
	return nil
}

// 压缩单个validator的收益记录
func (s *stakingRewardService) CompressValidatorReward(ctx context.Context, validator string, dayTime int64) error {
	endTime := dayTime + DAY_TIME
	powerSnapshots, err := s.GetPowerSnapshotsByValidator(ctx, validator, dayTime, endTime)
	if err != nil {
		return err
	}
	if len(powerSnapshots) == 0 {
		return nil
	}
	validatorRewardRecords, err := s.ListValidatorReward(ctx, validator, dayTime)
	if err != nil {
		return err
	}
	if len(validatorRewardRecords) == 0 {
		return nil
	}
	dayTimeList := make([]*ValidatorReward, 0)

	psIndex := -1
	totalShare := 0.0
	var validatorReward *ValidatorReward = nil
	for _, record := range validatorRewardRecords {
		for psIndex < len(powerSnapshots)-1 && record.Timestamp > powerSnapshots[psIndex+1].Timestamp {
			psIndex++
			totalShare = powerSnapshots[psIndex].ValidatorTotalShares
			validatorReward = nil
		}
		if validatorReward == nil {
			validatorReward = &ValidatorReward{
				RewardAmount: record.RewardAmount - record.CommissionAmount,
				Timestamp:    record.Timestamp,
				TotalShare:   totalShare,
			}
			dayTimeList = append(dayTimeList, validatorReward)
		} else {
			validatorReward.RewardAmount += record.RewardAmount - record.CommissionAmount
			//validatorReward.TotalShare += totalShare
		}
	}
	err = s.SetValidatorRewardToCache(ctx, validator, dayTime, dayTimeList)
	if err != nil {
		return err
	}
	return nil
}
