package main

import (
	"app/internal/cache/redisclient"
	"app/internal/config"
	"app/internal/db/postgres"
	"app/internal/service"
	"app/pkg/utils"
	"app/pkg/utils/otelutil"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	config.Parse("")
	redisclient.Load()
	postgres.Load()
	// 创建带取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	otelOptions := config.GetOtelOptions()
	otelOptions.ServiceName = "desert-staking-service"
	// Set up OpenTelemetry.
	otelShutdown, err := otelutil.SetupOTelSDK(ctx, otelOptions)
	if err != nil {
		log.Fatalf("failed to initialize otel sdk: %v", err)
	}

	// 启动定时任务
	cron := startCron(ctx)

	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigChan
	utils.Logger.Info(ctx, "main", "shutdown", "received shutdown signal")
	shutdown(ctx, cron, otelShutdown)

}

// 异步优雅关闭服务
func shutdown(ctx context.Context, cron *cron.Cron, otelShutdown func(context.Context) error) {
	// 创建带超时的上下文用于优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// 异步执行关闭操作
	done := make(chan struct{})
	go func() {
		release(ctx, cron, otelShutdown)
		close(done)
	}()

	// 等待关闭完成或超时
	select {
	case <-shutdownCtx.Done():
		utils.Logger.Error(ctx, "main", "shutdown", "shutdown timeout")
		os.Exit(1)
	case <-done:
		utils.Logger.Info(ctx, "main", "shutdown", "graceful shutdown completed")
	}

}

// 释放资源
func release(ctx context.Context, cron *cron.Cron, otelShutdown func(context.Context) error) {
	// 停止定时任务
	cron.Stop()
	utils.Logger.Info(ctx, "main", "shutdown", "cron stopped")

	// 关闭数据库连接
	if sqlDB, err := postgres.Default().DB(); err != nil {
		utils.Logger.Error(ctx, "main", "shutdown", "failed to get postgres db instance", slog.Attr{
			Key:   "err",
			Value: slog.AnyValue(err),
		})
	} else {
		if err := sqlDB.Close(); err != nil {
			utils.Logger.Error(ctx, "main", "shutdown", "failed to close postgres connection", slog.Attr{
				Key:   "err",
				Value: slog.AnyValue(err),
			})
		}
	}

	// 关闭 Redis 连接
	if err := redisclient.RedisClient.Close(); err != nil {
		utils.Logger.Error(ctx, "main", "shutdown", "failed to close redis connection", slog.Attr{
			Key:   "err",
			Value: slog.AnyValue(err),
		})
	}
	// 关闭 OpenTelemetry
	if err := otelShutdown(context.Background()); err != nil {
		utils.Logger.Error(ctx, "main", "shutdown", "failed to close otel connection", slog.Attr{
			Key:   "err",
			Value: slog.AnyValue(err),
		})
	}
}

// 定时任务
func startCron(ctx context.Context) *cron.Cron {
	// c := cron.New(cron.WithSeconds())
	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithSeconds())

	//DataAxgradService
	c.AddFunc("*/1 * * * * *", func() {
		err := service.ValidatorDescriptionService.SyncValidatorDescriptionAvatar(ctx)
		if err != nil {
			utils.Logger.Error(ctx, "validator_description", "SyncValidatorDescriptionAvatar", "sync validator description avatar", slog.Attr{
				Key:   "err",
				Value: slog.AnyValue(err),
			})
		}
	})
	c.AddFunc("*/1 * * * * *", func() {
		err := service.StakingRewardService.StartValidatorRewardCompressedTask(ctx)
		if err != nil {
			utils.Logger.Error(ctx, "staking_reward", "StartValidatorRewardCompressedTask", "start validator reward compressed task", slog.Attr{
				Key:   "err",
				Value: slog.AnyValue(err),
			})
		}
	})

	c.Start()
	utils.Logger.Info(ctx, "cron", "startCron", "cron start")
	return c
}
