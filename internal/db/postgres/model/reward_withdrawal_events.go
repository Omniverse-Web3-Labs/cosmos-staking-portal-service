package model

const RewardWithdrawalEventsTableName = "reward_withdrawal_events"

type RewardWithdrawalEvents struct {
	ID          string  `json:"id"`
	BlockHeight int     `json:"block_height"`
	TxHash      string  `json:"tx_hash"`
	Validator   string  `json:"validator"`
	Delegator   string  `json:"delegator"`
	Fee         float64 `json:"fee"`
	Amount      float64 `json:"amount"`
	Timestamp   int64   `json:"timestamp"`
}

func (r *RewardWithdrawalEvents) TableName() string {
	return RewardWithdrawalEventsTableName
}
