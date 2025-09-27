package model

import "app/pkg/utils"

type ValidatorEventQuery struct {
	utils.PageRequest
	Validator string `json:"validator"`
}

type ValidatorEvent struct {
	ID                 string  `json:"id"`
	Validator          string  `json:"validator"`
	TxHash             string  `json:"tx_hash"`
	BlockHeight        int     `json:"block_height"`
	EventType          string  `json:"event_type"`
	PrevCommissionRate float64 `json:"prev_commission_rate"`
	CommissionRate     float64 `json:"commission_rate"`
	MinSelfDelegation  float64 `json:"min_self_delegation"`
	Timestamp          int64   `json:"timestamp"`
	Age                int64   `json:"age"`
}

type PowerEventQuery struct {
	utils.PageRequest
	Validator string `json:"validator"`
	Delegator string `json:"delegator"`
}

type PowerEvent struct {
	ID          string  `json:"id"`
	BlockHeight int     `json:"block_height"`
	TxHash      string  `json:"tx_hash"`
	Validator   string  `json:"validator"`
	Delegator   string  `json:"delegator"`
	Amount      float64 `json:"amount"`
	EventType   string  `json:"event_type"`
	Timestamp   int64   `json:"timestamp"`
	Age         int64   `json:"age"`
}
