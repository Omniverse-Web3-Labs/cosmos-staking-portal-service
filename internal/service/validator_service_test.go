package service

import (
	"context"
	"fmt"
	"testing"
)

func TestGetValidatorListSortedByAddress(t *testing.T) {
	fmt.Println("validator list sorted by address, asc")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "address", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("address", v.Address)
	}

	fmt.Println("validator list sorted by address, desc")
	validators, _, err = ValidatorService.GetValidatorList(context.Background(), "address", true, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("address", v.Address)
	}
}

func TestGetValidatorListSortedByName(t *testing.T) {
	fmt.Println("validator list sorted by name, asc")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "name", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("name", v.Name)
	}

	fmt.Println("validator list sorted by name, desc")
	validators, _, err = ValidatorService.GetValidatorList(context.Background(), "name", true, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("name", v.Name)
	}
}

func TestGetValidatorListSortedByStake(t *testing.T) {
	fmt.Println("validator list sorted by stake, asc")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "stake", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("stake", v.Stake)
	}

	fmt.Println("validator list sorted by stake, desc")
	validators, _, err = ValidatorService.GetValidatorList(context.Background(), "stake", true, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("stake", v.Stake)
	}
}

func TestGetValidatorListSortedByVote(t *testing.T) {
	fmt.Println("validator list sorted by vote, asc")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "vote", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("vote", v.Voting)
	}

	fmt.Println("validator list sorted by vote, desc")
	validators, _, err = ValidatorService.GetValidatorList(context.Background(), "vote", true, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("vote", v.Voting)
	}
}

func TestGetValidatorListSortedByCommission(t *testing.T) {
	fmt.Println("validator list sorted by commission, asc")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "commission", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("commission", v.Commission)
	}

	fmt.Println("validator list sorted by commission, desc")
	validators, _, err = ValidatorService.GetValidatorList(context.Background(), "commission", true, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range validators {
		fmt.Println("commission", v.Commission)
	}
}

func TestGetDelegators(t *testing.T) {
	fmt.Println("get delegators")
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "commission", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	delegators, _, err := ValidatorService.GetDelegators(context.Background(), validators[0].Address, "", 0, 10, true)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	for _, v := range delegators {
		fmt.Printf("delegator: %+v\n", v)
	}
}

func TestGetValidators(t *testing.T) {
	fmt.Println("get validators")
	validatorList, _, err := ValidatorService.GetValidatorList(context.Background(), "commission", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	validatorAddresses := make([]string, 0)
	for _, v := range validatorList {
		validatorAddresses = append(validatorAddresses, v.Address)
	}

	validators, err := ValidatorService.GetValidators(context.Background(), validatorAddresses)
	if err != nil {
		t.Fatal("error", err.Error())
	}
	fmt.Println("validators", len(validators))
}

func TestGetStakeOverview(t *testing.T) {
	overview, err := ValidatorService.GetStakeOverview(context.Background())
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("overview: %+v\n", overview)
}

func TestGetValidatorDetail(t *testing.T) {
	validators, _, err := ValidatorService.GetValidatorList(context.Background(), "commission", false, -1, 0, "")
	if err != nil {
		t.Fatal("error", err.Error())
	}

	validator, err := ValidatorService.GetValidatorDetail(context.Background(), validators[0].Address)
	if err != nil {
		t.Fatal("error", err.Error())
	}

	fmt.Printf("validator: %+v\n", validator)
}
