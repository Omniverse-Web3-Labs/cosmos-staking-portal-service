package service

import (
	"app/internal/cache/redisclient"
	"app/internal/config"
	"context"
	"fmt"
	"testing"
)

func TestGetRawValidators(t *testing.T) {
	config.Parse("../../configs/config.yaml")
	redisclient.Load()
	validators, err := ChainDataService.GetValidators(context.Background())
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Println("update time", validators[0].Commission.UpdateTime)
}

func TestGetRawDelegators(t *testing.T) {
	validators, err := ChainDataService.GetValidators(context.Background())
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegation, err := ChainDataService.GetDelegators(context.Background(), validators[0].OperatorAddress, "", 0, 10, true)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("raw delegators:%+v\n", delegation)
}

func TestGetValidator(t *testing.T) {
	validators, err := ChainDataService.GetValidators(context.Background())
	if err != nil {
		t.Fatal("error", err.Error())
	}

	validator, err := ChainDataService.GetValidator(context.Background(), validators[0].OperatorAddress)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("validator:%+v", validator)
}

func TestGetDelegation(t *testing.T) {
	validators, err := ChainDataService.GetValidators(context.Background())
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegation, err := ChainDataService.GetDelegators(context.Background(), validators[0].OperatorAddress, "", 0, 10, true)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("raw delegation:%+v\n", delegation)
}

func TestGetPool(t *testing.T) {
	_, err := ChainDataService.GetPool(context.Background(), false)
	if err != nil {
		t.Fatal("error", err.Error())
	}
}

func TestGetSupply(t *testing.T) {
	_, err := ChainDataService.GetTotalSupply(context.Background(), false)
	if err != nil {
		t.Fatal("error", err.Error())
	}
}

func TestGetBlock(t *testing.T) {
	block, err := ChainDataService.GetBlock(context.Background(), 1)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("block: %+v\n", block)
}
