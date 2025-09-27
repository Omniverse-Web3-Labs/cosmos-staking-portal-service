package model

const TableNameStakingRewardDaily = "staking_reward_daily"

type StakingRewardDaily struct {
	ID            int64   `json:"id"`
	Delegator     string  `json:"delegator"`
	RewardAmount  float64 `json:"reward_amount"`
	StakingAmount float64 `json:"staking_amount"`
	DayTime       int64   `json:"day_time"`
	CreatedAt     int64   `json:"created_at"`
}

func (StakingRewardDaily) TableName() string {
	return TableNameStakingRewardDaily
}
