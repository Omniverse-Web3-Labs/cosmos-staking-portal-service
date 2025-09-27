package warmup

import (
	"context"
	"log"
	"time"
)

const TTL_COUNTDOWN_ALL = 60 * time.Second

func StartPrewarm(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = warmupAvatar(ctx)
			_ = warmupNetwork(ctx)
		case <-ctx.Done():
			log.Println("warmup shutdown.")
			return
		}
	}
}
