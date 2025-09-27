package warmup

import (
	"app/internal/cache/redisclient"
	"app/internal/service"
	"context"
	"fmt"
	"regexp"
)

func warmupAvatar(ctx context.Context) error {
	rawValidators, err := service.ChainDataService.GetValidators(ctx)
	if err != nil {
		return err
	}

	var identityRe = regexp.MustCompile(`^[a-fA-F0-9]{16}$`)
	for _, v := range rawValidators {
		identity := v.Description.Identity

		matched := identityRe.MatchString(identity)

		if !matched {
			continue
		}

		redisKey := fmt.Sprintf(redisclient.KeyValidatorAvatar, identity)

		ttl, err := redisclient.RedisClient.TTL(ctx, redisKey).Result()
		if err != nil {
			fmt.Printf("get redis ttl %s failed with err: %s\n", redisKey, err.Error())
			continue
		}

		if ttl < TTL_COUNTDOWN_ALL {
			_, err := service.ChainDataService.GetValidatorAvatar(ctx, identity, true)
			if err != nil {
				fmt.Printf("get validator avatar %s failed with err: %s\n", identity, err.Error())
				continue
			}
		}
	}

	return nil
}
