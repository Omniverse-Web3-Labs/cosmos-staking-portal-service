package model

const PowerEventsTableName = "power_events"

type PowerEvents struct {
	ID          string  `json:"id"`
	BlockHeight int     `json:"block_height"`
	TxHash      string  `json:"tx_hash"`
	Validator   string  `json:"validator"`
	Delegator   string  `json:"delegator"`
	Amount      float64 `json:"amount"`
	EventType   string  `json:"event_type"`
	Timestamp   int64   `json:"timestamp"`
}

func (p *PowerEvents) TableName() string {
	return PowerEventsTableName
}
