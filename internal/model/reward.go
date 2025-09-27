package model

import (
	"app/pkg/utils"
	"encoding/json"
)

// Delegator每日奖励
type StakingRewardDailyQuery struct {
	Delegator string `json:"delegator"`
	// 7: 7d, 30: 30d, 90: 90d
	Range json.Number `json:"range"`
}

type StakingRewardDaily struct {
	Timestamp     int64   `json:"timestamp"`
	RewardAmount  float64 `json:"reward"`
	StakingAmount float64 `json:"staking"`
}

// Delegator收益率
type StakingRewardRateQuery struct {
	Delegator string `json:"delegator"`
	// 1: 1d, 7: 7d, 30: 30d
	Period json.Number `json:"period"`
}

type StakingRewardRate struct {
	Timestamp int64   `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

type StakingRewardWithdrawalQuery struct {
	utils.PageRequest
	Delegator string `json:"delegator"`
}

type RewardWithdrawalEvents struct {
	ID          string  `json:"id"`
	BlockHeight int     `json:"block_height"`
	TxHash      string  `json:"tx_hash"`
	Validator   string  `json:"validator"`
	Delegator   string  `json:"delegator"`
	Fee         float64 `json:"fee"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
	Moniker     string  `json:"moniker"`
	Picture     string  `json:"picture"`
}
