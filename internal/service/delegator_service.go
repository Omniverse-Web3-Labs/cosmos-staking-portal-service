package service

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

type Delegation struct {
	Name       string                `json:"name"`
	Validator  string                `json:"validator"`
	Stake      string                `json:"stake"`      // big int with decimal 18, 1 OMNI: "1000000000000000000"
	Commission int                   `json:"commission"` // 0.01%. 1234 -> 12.34%
	Active     bool                  `json:"active"`
	Uptime     int                   `json:"uptime"`
	Picture    string                `json:"picture"`
	Amount     string                `json:"amount"`
	APY        string                `json:"apy"` // 0.01%. 1234 -> 12.34%
	Unbondings []UnbondingDelegation `json:"unbondings"`
}

type UnbondingDelegation struct {
	Amount    string    `json:"amount"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}

type DelegatorDetail struct {
	TotalStake string `json:"stake"`
	Earned     string `json:"earned"` // float64
	APY        string `json:"apy"`
}

type StakingParams struct {
	UnbondingTime float64 `json:"unbondingTime"` // days
	MaxEntries    int64   `json:"maxEntries"`
}

var DelegatorService delegatorService

type delegatorService struct {
}

func (service *delegatorService) GetUserDelegations(ctx context.Context, delegator string) ([]*Delegation, error) {
	rawDelegators, err := ChainDataService.GetDelegations(ctx, delegator)
	if err != nil {
		return nil, err
	}

	rawUnbondingDelegations, err := ChainDataService.GetUnbondingDelegations(ctx, delegator)
	if err != nil {
		return nil, err
	}

	avgAPY, err := ValidatorService.GetStakeAPY(ctx)
	if err != nil {
		return nil, err
	}

	validatorAddresses := make([]string, 0)
	delegations := make([]*Delegation, 0)
	for _, d := range rawDelegators {
		rawValidator, err := ChainDataService.GetValidator(ctx, d.Delegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		commission := float64(rawValidator.Commission.CommissionRates.Rate)

		active := true
		if rawValidator.Status != "BOND_STATUS_BONDED" {
			active = false
		}

		picture, err := ChainDataService.GetValidatorAvatar(ctx, rawValidator.Description.Identity, false)

		if err != nil {
			fmt.Println("Get validator avatar failed", err.Error())
		}

		validator := ValidatorDetail{
			Name:       rawValidator.Description.Moniker,
			Address:    rawValidator.OperatorAddress,
			Stake:      rawValidator.Tokens.Int.String(),
			Commission: int(commission * 10000),
			Identity:   rawValidator.Description.Identity,
			Picture:    picture,
			Active:     active,
		}

		rate := big.NewFloat(0).Sub(big.NewFloat(1), big.NewFloat(float64(rawValidator.Commission.CommissionRates.Rate)))

		delegation := Delegation{
			Name:       validator.Name,
			Validator:  validator.Address,
			Stake:      validator.Stake,
			Commission: validator.Commission,
			Active:     validator.Active,
			Uptime:     0,
			Picture:    validator.Picture,
			Amount:     d.Balance.Amount.String(),
			APY:        big.NewFloat(0).Mul(big.NewFloat(0).Mul(avgAPY, rate), big.NewFloat(10000)).String(),
			Unbondings: make([]UnbondingDelegation, 0),
		}
		validatorAddresses = append(validatorAddresses, validator.Address)
		delegations = append(delegations, &delegation)
	}

	findDelegation := func(delegations []*Delegation, validator string) *Delegation {
		for _, d := range delegations {
			if d.Validator == validator {
				return d
			}
		}

		return nil
	}

	for _, unbonding := range rawUnbondingDelegations {
		d := findDelegation(delegations, unbonding.ValidatorAddress)
		if d == nil {
			rawValidator, err := ChainDataService.GetValidator(ctx, unbonding.ValidatorAddress)
			if err != nil {
				return nil, err
			}

			commission := float64(rawValidator.Commission.CommissionRates.Rate)

			active := true
			if rawValidator.Status != "BOND_STATUS_BONDED" {
				active = false
			}

			picture, err := ChainDataService.GetValidatorAvatar(ctx, rawValidator.Description.Identity, false)

			if err != nil {
				fmt.Println("Get validator avatar failed", err.Error())
			}

			validator := ValidatorDetail{
				Name:       rawValidator.Description.Moniker,
				Address:    rawValidator.OperatorAddress,
				Stake:      rawValidator.Tokens.Int.String(),
				Commission: int(commission * 10000),
				Identity:   rawValidator.Description.Identity,
				Picture:    picture,
				Active:     active,
			}

			rate := big.NewFloat(0).Sub(big.NewFloat(1), big.NewFloat(float64(rawValidator.Commission.CommissionRates.Rate)))

			d = &Delegation{
				Name:       validator.Name,
				Validator:  validator.Address,
				Stake:      validator.Stake,
				Commission: validator.Commission,
				Active:     validator.Active,
				Uptime:     0,
				Picture:    validator.Picture,
				Amount:     "0",
				APY:        big.NewFloat(0).Mul(big.NewFloat(0).Mul(avgAPY, rate), big.NewFloat(10000)).String(),
				Unbondings: make([]UnbondingDelegation, 0),
			}
			delegations = append(delegations, d)
			validatorAddresses = append(validatorAddresses, validator.Address)
		}

		for _, entry := range unbonding.Entries {
			rawBlock, err := ChainDataService.GetBlock(ctx, int64(entry.CreationHeight))
			if err != nil {
				return nil, err
			}
			d.Unbondings = append(d.Unbondings, UnbondingDelegation{
				Amount:    entry.Balance.String(),
				StartTime: rawBlock.Block.Header.Time,
				EndTime:   entry.CompletionTime,
			})
		}
	}

	// uptime
	slashingParams, err := ChainDataService.GetSlashingParams(ctx)
	if err != nil {
		return nil, err
	}

	blockWindow := int(slashingParams.Params.SignedBlocksWindow)
	uptimes, err := ValidatorUptimeService.ListRecentUptimeByOperAddresses(ctx, validatorAddresses, blockWindow)
	if err != nil {
		return nil, err
	}

	uptimeMap := make(map[string]int)
	for _, v := range uptimes {
		uptimeMap[v.Validator] = v.Count * 10000 / blockWindow
	}

	for _, v := range delegations {
		v.Uptime = uptimeMap[v.Validator]
	}

	return delegations, nil
}

func (service *delegatorService) GetDelegatorDetail(ctx context.Context, delegator string) (*DelegatorDetail, error) {
	rawDelegators, err := ChainDataService.GetDelegations(ctx, delegator)
	if err != nil {
		return nil, err
	}

	totalStake := big.NewInt(0)
	rewardPerYear := big.NewFloat(0)
	avgAPY, err := ValidatorService.GetStakeAPY(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range rawDelegators {
		totalStake.Add(totalStake, d.Balance.Amount.Int)
		rawValidator, err := ChainDataService.GetValidator(ctx, d.Delegation.ValidatorAddress)
		if err != nil {
			return nil, err
		}

		rate := big.NewFloat(0).Sub(big.NewFloat(1), big.NewFloat(float64(rawValidator.Commission.CommissionRates.Rate)))
		reward := big.NewFloat(0).Mul(big.NewFloat(0).Mul(rate, big.NewFloat(0).SetInt(d.Balance.Amount.Int)), avgAPY)
		rewardPerYear.Add(rewardPerYear, reward)
	}

	totalStakeFloat := big.NewFloat(0).SetInt(totalStake)
	apy := big.NewFloat(0)
	if totalStakeFloat.Cmp(big.NewFloat(0)) != 0 {
		apy.Mul(big.NewFloat(0).Quo(rewardPerYear, totalStakeFloat), big.NewFloat(10000))
	}
	apyInt := big.NewInt(0)
	apy.Int(apyInt)

	earnedReward, err := StakingRewardService.GetDelegatorTotalReward(ctx, delegator)
	if err != nil {
		return nil, err
	}

	delegatorDetail := DelegatorDetail{
		TotalStake: totalStake.String(),
		Earned:     fmt.Sprintf("%f", earnedReward),
		APY:        apyInt.String(),
	}

	return &delegatorDetail, nil
}

func (service *delegatorService) GetStakingParams(ctx context.Context) (*StakingParams, error) {
	params, err := ChainDataService.GetStakingParams(ctx)
	if err != nil {
		return nil, err
	}

	unbondingTime, err := time.ParseDuration(params.Params.UnbondingTime)
	if err != nil {
		return nil, err
	}

	return &StakingParams{
		UnbondingTime: unbondingTime.Hours() / 24,
		MaxEntries:    params.Params.MaxEntries,
	}, nil
}
