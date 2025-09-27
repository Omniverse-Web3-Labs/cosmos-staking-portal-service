package service

import (
	"context"
	"fmt"
	"testing"
)

func TestGetUserDelegations(t *testing.T) {
	fmt.Println("get user delegations")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "voting", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegators, _, err := ValidatorService.GetDelegators(context.Background(), validators[0].Address, "", 0, 10, true)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegations, err := DelegatorService.GetUserDelegations(context.Background(), delegators[0].Address)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Println("delegations:")
	for _, d := range delegations {
		fmt.Printf(" %+v\n", *d)
	}
}

func TestGetDelegatorDetail(t *testing.T) {
	fmt.Println("get delegator detail")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "voting", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegators, _, err := ValidatorService.GetDelegators(context.Background(), validators[0].Address, "", 0, 10, true)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegations, err := DelegatorService.GetDelegatorDetail(context.Background(), delegators[0].Address)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("delegator detail: %+v\n", delegations)
}
