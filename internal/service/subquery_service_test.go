package service

import (
	"app/internal/cache/redisclient"
	"app/internal/config"
	"app/internal/db/postgres"
	"app/internal/model"
	"app/pkg/utils"
	"context"
	"fmt"
	"os"
	"testing"
)

func setup() {
	fmt.Println("全局 setup")
	// 初始化代码，如建立数据库连接、加载配置等
	config.Parse("../../configs/config.yaml")
	postgres.Load()
	redisclient.Load()
}

func teardown() {
	fmt.Println("全局 teardown")
	// 清理资源
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run() // 执行所有测试
	teardown()
	os.Exit(code)
}

func TestGetValidatorEvents(t *testing.T) {
	ctx := context.Background()
	params := &model.ValidatorEventQuery{
		Validator: "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg",
		PageRequest: utils.PageRequest{
			Current: 3,
			Size:    1,
			Order:   "asc",
			Sort:    "block_height",
		},
	}
	events, err := SubqueryService.GetValidatorEvents(ctx, params)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("raw events:%+v\n", events)
}

func TestGetPowerEvents(t *testing.T) {
	ctx := context.Background()
	params := &model.PowerEventQuery{
		Validator: "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg",
		Delegator: "omnivaloper195tf7tuvxagrsq7seh8579qrjvxc6j9wws29sg",
		PageRequest: utils.PageRequest{
			Current: 1,
			Size:    1,
			Order:   "asc",
			Sort:    "block_height",
		},
	}
	events, err := SubqueryService.GetPowerEvents(ctx, params)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("raw events:%+v\n", events)
}

func TestGetUnbondingDelegations(t *testing.T) {
	validators, err := ChainDataService.GetValidators(context.Background())
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegators, err := ChainDataService.GetDelegators(context.Background(), validators[0].OperatorAddress, "", 0, 10, true)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegation, err := ChainDataService.GetUnbondingDelegations(context.Background(), delegators[0].Delegation.DelegatorAddress)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("raw unbonding delegation:%+v\n", delegation)
}
