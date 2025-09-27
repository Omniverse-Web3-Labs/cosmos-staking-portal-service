package service

import (
	"app/internal/db/postgres"
	dbModel "app/internal/db/postgres/model"
	"app/internal/model"
	"app/pkg/utils"
	"context"
	"time"
)

var SubqueryService subqueryService

type subqueryService struct {
}

// historical data
func (service *subqueryService) GetValidatorEvents(ctx context.Context, query *model.ValidatorEventQuery) (*utils.Paginator, error) {
	db := postgres.Subquery()
	paginator := utils.NewPaginator(query.Current, query.Size)
	list := make([]*model.ValidatorEvent, 0)
	err := db.Model(&dbModel.ValidatorEvents{}).WithContext(ctx).
		Where("validator = ?", query.Validator).
		Select("block_height, tx_hash, validator, event_type, prev_commission_rate, commission_rate, min_self_delegation, timestamp").
		Order("timestamp DESC").Limit(paginator.Size + 1).Offset(paginator.Offset()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	paginator.HasMore = len(list) > paginator.Size
	if paginator.HasMore {
		list = list[:len(list)-1]
	}
	now := time.Now().UnixMilli()
	for _, v := range list {
		v.Age = now - v.Timestamp
	}
	paginator.SetList(list)
	return paginator, nil
}

// historical data
func (service *subqueryService) GetPowerEvents(ctx context.Context, query *model.PowerEventQuery) (*utils.Paginator, error) {
	db := postgres.Subquery()
	paginator := utils.NewPaginator(query.Current, query.Size)
	list := make([]*model.PowerEvent, 0)
	err := db.Model(&dbModel.PowerEvents{}).WithContext(ctx).
		Where("validator = ?", query.Validator).
		Select("block_height, tx_hash, validator, delegator, amount, event_type, timestamp").
		Order("timestamp DESC").Limit(paginator.Size + 1).Offset(paginator.Offset()).Find(&list).Error
	if err != nil {
		return nil, err
	}
	paginator.HasMore = len(list) > paginator.Size
	if paginator.HasMore {
		list = list[:len(list)-1]
	}
	now := time.Now().UnixMilli()
	for _, v := range list {
		v.Age = now - v.Timestamp
	}
	paginator.SetList(list)
	return paginator, nil
}
