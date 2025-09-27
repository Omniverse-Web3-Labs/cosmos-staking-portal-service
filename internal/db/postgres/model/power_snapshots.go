package model

// CREATE TABLE "public"."power_snapshots" (
// 	"id" text COLLATE "pg_catalog"."default" NOT NULL,
// 	"block_height" int4 NOT NULL,
// 	"validator" text COLLATE "pg_catalog"."default" NOT NULL,
// 	"delegator" text COLLATE "pg_catalog"."default" NOT NULL,
// 	"amount" float8 NOT NULL,
// 	"delegator_total_shares" float8 NOT NULL,
// 	"validator_total_shares" float8 NOT NULL,
// 	"validator_total_amount" float8 NOT NULL,
// 	"timestamp" numeric NOT NULL,
// 	"created_at" numeric NOT NULL,
// 	"_id" uuid NOT NULL,
// 	"_block_range" int8range NOT NULL,
// 	CONSTRAINT "power_snapshots_pkey" PRIMARY KEY ("_id")
//   )
//   ;
type PowerSnapshots struct {
	ID                   string  `gorm:"column:id;primaryKey" json:"id"`
	Validator            string  `gorm:"column:validator" json:"validator"`
	Delegator            string  `gorm:"column:delegator" json:"delegator"`
	BlockHeight          int     `gorm:"column:block_height" json:"block_height"`
	Amount               float64 `gorm:"column:amount" json:"amount"`
	DelegatorTotalShares float64 `gorm:"column:delegator_total_shares" json:"delegator_total_shares"`
	ValidatorTotalShares float64 `gorm:"column:validator_total_shares" json:"validator_total_shares"`
	ValidatorTotalAmount float64 `gorm:"column:validator_total_amount" json:"validator_total_amount"`
	Timestamp            int64   `gorm:"column:timestamp" json:"timestamp"`
	CreatedAt            int64   `gorm:"column:created_at" json:"created_at"`
}
