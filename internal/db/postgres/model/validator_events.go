package model

const ValidatorEventsTableName = "validator_events"

type ValidatorEvents struct {
	ID                 string  `json:"id"`
	Validator          string  `json:"validator"`
	TxHash             string  `json:"tx_hash"`
	BlockHeight        int     `json:"block_height"`
	EventType          string  `json:"event_type"`
	PrevCommissionRate float64 `json:"prev_commission_rate"`
	CommissionRate     float64 `json:"commission_rate"`
	MinSelfDelegation  int64   `json:"min_self_delegation"`
	Timestamp          int64   `json:"timestamp"`
}

func (v *ValidatorEvents) TableName() string {
	return ValidatorEventsTableName
}
