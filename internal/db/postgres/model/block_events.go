package model

const BlockEventsTableName = "uptime.block_events"

type BlockEvents struct {
	ID              string `json:"id"`
	BlockHeight     int    `json:"block_height"`
	ProposerAddress string `json:"proposer_address"`
	Timestamp       int64  `json:"timestamp"`
}

func (p *BlockEvents) TableName() string {
	return BlockEventsTableName
}
