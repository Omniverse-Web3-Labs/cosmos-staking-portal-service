package warmup

import (
	"app/internal/cache/redisclient"
	"app/internal/service"
	"context"
	"fmt"
	"time"
)

const TTL_COUNTDOWN_STAKING_POOL = 1 * time.Minute
const TTL_COUNTDOWN_TOTAL_SUPPLY = 1 * time.Minute
const TTL_COUNTDOWN_INFLATION = 1 * time.Minute

func warmupNetwork(ctx context.Context) error {
	var err error
	err = warmupStakingPool(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = warmupTotalSupply(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = warmupInflation(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func warmupStakingPool(ctx context.Context) error {
	ttl, err := redisclient.RedisClient.TTL(ctx, redisclient.KeyRawPool).Result()
	if err != nil {
		return fmt.Errorf("get redis ttl %s failed with err: %s", redisclient.KeyRawPool, err.Error())
	}

	if ttl < TTL_COUNTDOWN_ALL {
		_, err := service.ChainDataService.GetPool(ctx, true)
		if err != nil {
			return fmt.Errorf("get staking pool failed with err: %s", err.Error())
		}
	}

	return nil
}

func warmupTotalSupply(ctx context.Context) error {
	ttl, err := redisclient.RedisClient.TTL(ctx, redisclient.KeyRawSupply).Result()
	if err != nil {
		return fmt.Errorf("get redis ttl %s failed with err: %s", redisclient.KeyRawSupply, err.Error())
	}

	if ttl < TTL_COUNTDOWN_ALL {
		_, err := service.ChainDataService.GetTotalSupply(ctx, true)
		if err != nil {
			return fmt.Errorf("get total supply failed with err: %s", err.Error())
		}
	}

	return nil
}

func warmupInflation(ctx context.Context) error {
	ttl, err := redisclient.RedisClient.TTL(ctx, redisclient.KeyRawInflation).Result()
	if err != nil {
		return fmt.Errorf("get redis ttl %s failed with err: %s", redisclient.KeyRawInflation, err.Error())
	}

	if ttl < TTL_COUNTDOWN_ALL {
		_, err := service.ChainDataService.GetInflation(ctx, true)
		if err != nil {
			return fmt.Errorf("get inflation failed with err: %s", err.Error())
		}
	}

	return nil
}
