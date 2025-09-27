package main

import (
	"app/cmd/warmup/warmup"
	"app/internal/cache/redisclient"
	"app/internal/config"
	"context"
)

func main() {
	config.Parse("")
	redisclient.Load()

	warmup.StartPrewarm(context.Background())
}
