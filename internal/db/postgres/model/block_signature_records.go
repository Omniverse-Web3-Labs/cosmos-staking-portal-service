package model

const BlockSignatureRecordsTableName = "uptime.block_signature_records"

type BlockSignatureRecords struct {
	ID          string `json:"id"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	Timestamp   int64  `json:"timestamp"`
	Validator   string `json:"validator"`
	Signed      bool   `json:"signed"`
}

func (p *BlockSignatureRecords) TableName() string {
	return BlockSignatureRecordsTableName
}
