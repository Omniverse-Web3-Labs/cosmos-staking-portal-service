package service

import (
	"app/internal/cache/redisclient"
	"app/internal/db/postgres"
	dbModel "app/internal/db/postgres/model"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type validatorDescriptionService struct {
}

var ValidatorDescriptionService validatorDescriptionService

func (s *validatorDescriptionService) UpdateValidatorDescriptionAvatar(ctx context.Context, id string, avatar string) error {
	db := postgres.Subquery()
	err := db.Model(&dbModel.ValidatorDescriptions{}).WithContext(ctx).Where("id = ?", id).Update("avatar", avatar).Error
	if err != nil {
		return err
	}
	return nil
}

// 同步validatorDescription avatar 任务
func (s *validatorDescriptionService) SyncValidatorDescriptionAvatar(ctx context.Context) error {
	db := postgres.Subquery()
	list := make([]*dbModel.ValidatorDescriptions, 0)
	err := db.Model(&dbModel.ValidatorDescriptions{}).WithContext(ctx).Where("avatar = '' ").Find(&list).Error
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`^[a-fA-F0-9]{16}$`)
	for _, item := range list {
		var avatar string
		var err error
		var identity string = item.Identity
		matched := re.MatchString(identity)
		if !matched {
			avatar = "-"
			identity = "8F9F613D746F994A"
		}
		avatar, err = ChainDataService.GetValidatorAvatar(ctx, identity, false)
		if err != nil {
			return err
		}
		err = ValidatorDescriptionService.UpdateValidatorDescriptionAvatar(ctx, item.ID, avatar)
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取validator描述
func (s *validatorDescriptionService) GetValidatorDescription(ctx context.Context, validator string) (*dbModel.ValidatorDescriptions, error) {
	cacheKey := fmt.Sprintf("validator_description_%s", validator)
	cacheValue := redisclient.RedisClient.Get(ctx, cacheKey)
	if err := cacheValue.Err(); err != nil {
		if err == redis.Nil { // 缓存不存在,从数据库获取
			db := postgres.Subquery()
			validatorDescription := &dbModel.ValidatorDescriptions{}
			err = db.Model(&dbModel.ValidatorDescriptions{}).WithContext(ctx).Where("id = ?", validator).Order("timestamp DESC").First(validatorDescription).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			if err == gorm.ErrRecordNotFound {
				validatorDescription = &dbModel.ValidatorDescriptions{
					ID: validator,
				}
			}
			cacheValue, err := json.Marshal(validatorDescription)
			if err != nil {
				return nil, err
			}
			err = redisclient.RedisClient.Set(ctx, cacheKey, cacheValue, time.Minute*10).Err()
			if err != nil {
				return nil, err
			}
			return validatorDescription, nil
		}
		return nil, err
	}

	validatorDescription := &dbModel.ValidatorDescriptions{}
	cacheValueBytes, _ := cacheValue.Bytes()
	err := json.Unmarshal(cacheValueBytes, validatorDescription)
	if err != nil {
		return nil, err
	}
	return validatorDescription, nil
}
