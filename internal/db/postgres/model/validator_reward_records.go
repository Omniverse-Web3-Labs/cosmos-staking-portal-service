package model

type ValidatorRewardRecords struct {
	ID               string  `gorm:"column:id;primaryKey" json:"id"`
	BlockHeight      int     `gorm:"column:block_height" json:"block_height"`
	Validator        string  `gorm:"column:validator" json:"validator"`
	RewardAmount     float64 `gorm:"column:reward_amount" json:"reward_amount"`
	CommissionAmount float64 `gorm:"column:commission_amount" json:"commission_amount"`
	Timestamp        int64   `gorm:"column:timestamp" json:"timestamp"`
}
