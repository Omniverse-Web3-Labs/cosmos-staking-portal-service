package service

import (
	"app/internal/db/postgres"
	dbModel "app/internal/db/postgres/model"
	"context"
)

type blockEventService struct {
}

var BlockEventService blockEventService

// GetLatestAndPreviousBlock 获取最新区块和倒数第100个区块
func (s *blockEventService) GetLatestAndPreviousBlock(ctx context.Context) (*dbModel.BlockEvents, *dbModel.BlockEvents, error) {
	db := postgres.Subquery()

	var blocks []dbModel.BlockEvents
	err := db.Raw(`
		WITH ranked_blocks AS (
			SELECT *,
				ROW_NUMBER() OVER (ORDER BY timestamp DESC) as rn
			FROM uptime.block_events
		)
		SELECT *
		FROM ranked_blocks
		WHERE rn = 1 OR rn = 100
		ORDER BY timestamp DESC
	`).Find(&blocks).Error
	if err != nil {
		return nil, nil, err
	}

	if len(blocks) != 2 {
		return nil, nil, nil
	}

	return &blocks[0], &blocks[1], nil
}
