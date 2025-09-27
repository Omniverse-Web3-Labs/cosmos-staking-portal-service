package service

import (
	"app/pkg/utils"
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"
)

type Validator struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Stake      string `json:"stake"`      // big int with decimal 18, 1 OMNI: "1000000000000000000"
	Voting     int    `json:"voting"`     // 0.01%。 1234 -> 12.34%
	Commission int    `json:"commission"` // 0.01%. 1234 -> 12.34%
	APY        string `json:"apy"`        // 0.01%. 1234 -> 12.34%
	Picture    string `json:"picture"`
	Active     bool   `json:"active"`
	Uptime     int    `json:"uptime"` // 0.01%, 1234 -> 12.34%
}

type ValidatorDetail struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Stake      string `json:"stake"`      // big int with decimal 18, 1 OMNI: "1000000000000000000"
	Voting     int    `json:"voting"`     // 0.01%. 1234 -> 12.34%
	Commission int    `json:"commission"` // 0.01%. 1234 -> 12.34%
	Identity   string `json:"identity"`
	APY        string `json:"apy"` // 0.01%. 1234 -> 12.34%
	Picture    string `json:"picture"`
	Active     bool   `json:"active"`
	Uptime     int    `json:"uptime"` // 0.01%, 1234 -> 12.34%
}

type Delegator struct {
	Address                string `json:"address"`
	Stake                  string `json:"stake"` // big int with decimal 18, 1 OMNI: "1000000000000000000"
	StakingActivationEpoch int    `json:"activationEpoch"`
	HoldTime               int    `json:"holdTime"`
}

type CommissionRate struct {
	Date          time.Time `json:"date"`
	Previous      int       `json:"previous"`
	New           int       `json:"new"`
	Change        int       `json:"change"`
	EffectiveTime time.Time `json:"effectiveTime"`
}

type StakeOverview struct {
	TotalStake    string `json:"totalStake"`    // big int with decimal 18, 1 OMNI: "1000000000000000000"
	ActivelyStake string `json:"activelyStake"` // big int with decimal 18, 1 OMNI: "1000000000000000000"
	AvgAPY        string `json:"avgAPY"`        // 0.01%. 1234 -> 12.34%
}

var ValidatorService validatortService

type validatortService struct {
}

func (service *validatortService) GetValidatorList(ctx context.Context, sortBy string, desc bool, offset int, limit int, search string) ([]*Validator, int, error) {
	// if offset < 0 {
	// 	return nil, fmt.Errorf("offset should at least be 0")
	// }

	// if limit <= 0 {
	// 	return nil, fmt.Errorf("limit should at least be 1")
	// }

	// cache here
	validatorData, err := ChainDataService.GetValidators(ctx)
	if err != nil {
		return nil, 0, err
	}

	pool, err := ChainDataService.GetPool(ctx, false)
	if err != nil {
		return nil, 0, err
	}

	avgAPY, err := service.GetStakeAPY(ctx)
	if err != nil {
		return nil, 0, err
	}

	validators := make([]*Validator, 0)
	validatorAddresses := make([]string, 0)
	for _, v := range validatorData {
		stake := v.Tokens.Int
		votingPower := big.NewInt(0)
		if v.Status == "BOND_STATUS_BONDED" {
			votingPower = new(big.Int).Div(new(big.Int).Mul(stake, big.NewInt(10000)), pool.PoolData.BondedTokens.Int)
		}
		commission := float64(v.Commission.CommissionRates.Rate)
		rate := utils.NewBigFloatFromFloat64(1).Sub(big.NewFloat(commission))

		active := true
		if v.Status != "BOND_STATUS_BONDED" {
			active = false
		}

		picture, err := ChainDataService.GetValidatorAvatar(ctx, v.Description.Identity, false)

		if err != nil {
			fmt.Println("Get validator avatar failed", err.Error())
		}

		if search != "" {
			if strings.Contains(strings.ToLower(v.Description.Moniker), strings.ToLower(search)) || strings.Contains(strings.ToLower(v.OperatorAddress), strings.ToLower(search)) {
				validators = append(validators, &Validator{
					Name:       v.Description.Moniker,
					Address:    v.OperatorAddress,
					Stake:      v.Tokens.Int.String(),
					Voting:     int(votingPower.Int64()),
					Commission: int(commission * 10000),
					Picture:    picture,
					APY:        rate.Mul(avgAPY).Mul(big.NewFloat(10000)).Val().String(),
					Active:     active,
				})
				validatorAddresses = append(validatorAddresses, v.OperatorAddress)
			}
		} else {
			validators = append(validators, &Validator{
				Name:       v.Description.Moniker,
				Address:    v.OperatorAddress,
				Stake:      v.Tokens.Int.String(),
				Voting:     int(votingPower.Int64()),
				Commission: int(commission * 10000),
				Picture:    picture,
				APY:        rate.Mul(avgAPY).Mul(big.NewFloat(10000)).Val().String(),
				Active:     active,
			})
			validatorAddresses = append(validatorAddresses, v.OperatorAddress)
		}
	}

	// uptime
	slashingParams, err := ChainDataService.GetSlashingParams(ctx)
	if err != nil {
		return nil, 0, err
	}

	blockWindow := int(slashingParams.Params.SignedBlocksWindow)
	uptimes, err := ValidatorUptimeService.ListRecentUptimeByOperAddresses(ctx, validatorAddresses, blockWindow)
	if err != nil {
		return nil, 0, err
	}

	uptimeMap := make(map[string]int)
	for _, v := range uptimes {
		uptimeMap[v.Validator] = v.Count * 10000 / blockWindow
	}

	for _, v := range validators {
		v.Uptime = uptimeMap[v.Address]
	}

	// sort
	sort.Slice(validators, func(i, j int) bool {
		ret := false
		switch sortBy {
		case "address":
			ret = validators[i].Address < validators[j].Address
		case "name":
			ret = validators[i].Name < validators[j].Name
		case "stake":
			stakei := new(big.Int)
			stakei, ok := stakei.SetString(validators[i].Stake, 10)
			if !ok {
				fmt.Println("stake value error")
				return false
			}

			stakej := new(big.Int)
			stakej, ok = stakej.SetString(validators[j].Stake, 10)
			if !ok {
				fmt.Println("stake value error")
				return false
			}
			cmp := stakei.Cmp(stakej)
			if cmp == -1 {
				ret = true
			} else {
				ret = false
			}
		case "voting":
			ret = validators[i].Voting < validators[j].Voting
		case "commission":
			ret = validators[i].Commission < validators[j].Commission
		default:
			ret = validators[i].Voting < validators[j].Voting
		}
		if desc {
			ret = !ret
		}
		return ret
	})

	// // pagination
	// paging_validators := make([]Validator, 0)
	// for i := offset; i < offset+limit; i++ {
	// 	if len(validators) > i {
	// 		paging_validators = append(paging_validators, validators[i])
	// 	}
	// }
	return validators, len(validators), nil
}

func (service *validatortService) GetValidators(ctx context.Context, validatorAddresses []string) ([]*Validator, error) {
	pool, err := ChainDataService.GetPool(ctx, false)
	if err != nil {
		return nil, err
	}

	avgAPY, err := service.GetStakeAPY(ctx)
	if err != nil {
		return nil, err
	}

	validators := make([]*Validator, 0)
	for _, vaddr := range validatorAddresses {
		v, err := ChainDataService.GetValidator(ctx, vaddr)
		if err != nil {
			return nil, err
		}

		stake := v.Tokens.Int
		votingPower := big.NewInt(0)
		if v.Status == "BOND_STATUS_BONDED" {
			votingPower = new(big.Int).Div(new(big.Int).Mul(stake, big.NewInt(10000)), pool.PoolData.BondedTokens.Int)
		}
		commission := float64(v.Commission.CommissionRates.Rate)
		rate := utils.NewBigFloatFromFloat64(1).Sub(big.NewFloat(commission))

		active := true
		if v.Status != "BOND_STATUS_BONDED" {
			active = false
		}

		picture, err := ChainDataService.GetValidatorAvatar(ctx, v.Description.Identity, false)

		if err != nil {
			fmt.Println("Get validator avatar failed", err.Error())
		}

		validators = append(validators, &Validator{
			Name:       v.Description.Moniker,
			Address:    v.OperatorAddress,
			Stake:      v.Tokens.Int.String(),
			Voting:     int(votingPower.Int64()),
			Commission: int(commission * 10000),
			Picture:    picture,
			APY:        rate.Mul(avgAPY).Mul(big.NewFloat(10000)).Val().String(),
			Active:     active,
		})
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

	for _, v := range validators {
		v.Uptime = uptimeMap[v.Address]
	}

	return validators, nil
}

func (service *validatortService) GetDelegators(ctx context.Context, validatorAddress string, sortBy string, offset int, limit int, desc bool) ([]Delegator, int, error) {
	// cache here
	delegationData, err := ChainDataService.GetDelegators(ctx, validatorAddress, sortBy, offset, limit, desc)
	if err != nil {
		return nil, 0, err
	}

	totalShare := big.NewFloat(0)
	for _, v := range delegationData {
		totalShare.Add(totalShare, v.Delegation.Shares.Float)
	}

	delegators := make([]Delegator, 0)
	for _, v := range delegationData {
		// share := float64(v.Delegation.Shares)
		// stake := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(share/totalShare*1000000)), totalStake), big.NewInt(1000000)).String()
		delegators = append(delegators, Delegator{
			Address:                v.Delegation.DelegatorAddress,
			Stake:                  v.Balance.Amount.String(),
			StakingActivationEpoch: 0,
			HoldTime:               0,
		})
	}

	// pagination
	paging_delegators := make([]Delegator, 0)
	for i := offset; i < offset+limit; i++ {
		if len(delegators) > i {
			paging_delegators = append(paging_delegators, delegators[i])
		}
	}

	return paging_delegators, len(delegators), nil
}

func (service *validatortService) GetStakeOverview(ctx context.Context) (*StakeOverview, error) {
	// cache here

	bonded, err := ChainDataService.GetPool(ctx, false)
	if err != nil {
		return nil, err
	}

	apy, err := service.GetStakeAPY(ctx)
	if err != nil {
		return nil, err
	}

	validatorData, err := ChainDataService.GetValidators(ctx)
	if err != nil {
		return nil, err
	}

	weightReward := big.NewFloat(0)
	for _, v := range validatorData {
		rate := big.NewFloat(0).Sub(big.NewFloat(1), big.NewFloat(float64(v.Commission.CommissionRates.Rate)))
		weightReward.Add(weightReward, big.NewFloat(0).Mul(rate, big.NewFloat(0).SetInt(v.Tokens.Int)))
	}

	avgAPY := utils.NewBigFloatFromBigFloat(apy)
	if len(validatorData) > 0 {
		weightReward.Quo(weightReward, big.NewFloat(0).SetInt(bonded.PoolData.BondedTokens.Int))
		avgAPY.Mul(weightReward).Mul(big.NewFloat(10000))
	}

	stakeInfo := StakeOverview{
		TotalStake:    big.NewInt(0).Add(bonded.PoolData.BondedTokens.Int, bonded.PoolData.NotBondedTokens.Int).String(),
		ActivelyStake: bonded.PoolData.BondedTokens.String(),
		AvgAPY:        avgAPY.Val().String(),
	}

	return &stakeInfo, nil
}

func (service *validatortService) GetStakeAPY(ctx context.Context) (*big.Float, error) {
	supply, err := ChainDataService.GetTotalSupply(ctx, false)
	if err != nil {
		return nil, err
	}

	bonded, err := ChainDataService.GetPool(ctx, false)
	if err != nil {
		return nil, err
	}

	inflation, err := ChainDataService.GetInflation(ctx, false)
	if err != nil {
		return nil, err
	}

	distributionParams, err := ChainDataService.GetDistributionParams(ctx)
	if err != nil {
		return nil, err
	}
	tax := big.NewFloat(float64(distributionParams.Params.CommunityTax))
	proposerRewardRate := big.NewFloat(float64(distributionParams.Params.BaseProposerReward) + float64(distributionParams.Params.BonusProposerReward))

	infBigFloat := utils.NewBigFloatFromFloat64(inflation)
	avgAPY := infBigFloat.Mul(utils.Complement(tax)).Mul(utils.Complement(proposerRewardRate)).Mul(big.NewFloat(0).SetInt(supply.Amount.Int)).Div(big.NewFloat(0).SetInt(bonded.PoolData.BondedTokens.Int))

	return avgAPY.Val(), nil
}

func (service *validatortService) GetValidatorDetail(ctx context.Context, validatorAddress string) (*ValidatorDetail, error) {
	rawValidator, err := ChainDataService.GetValidator(ctx, validatorAddress)
	if err != nil {
		return nil, err
	}

	pool, err := ChainDataService.GetPool(ctx, false)
	if err != nil {
		return nil, err
	}

	stake := rawValidator.Tokens.Int
	votingPower := big.NewInt(0)
	if rawValidator.Status == "BOND_STATUS_BONDED" {
		votingPower = new(big.Int).Div(new(big.Int).Mul(stake, big.NewInt(10000)), pool.PoolData.BondedTokens.Int)
	}
	commission := float64(rawValidator.Commission.CommissionRates.Rate)

	active := true
	if rawValidator.Status != "BOND_STATUS_BONDED" {
		active = false
	}

	// uptime
	slashingParams, err := ChainDataService.GetSlashingParams(ctx)
	if err != nil {
		return nil, err
	}

	blockWindow := int(slashingParams.Params.SignedBlocksWindow)
	uptimes, err := ValidatorUptimeService.ListRecentUptimeByOperAddresses(ctx, []string{rawValidator.OperatorAddress}, blockWindow)
	if err != nil {
		return nil, err
	}

	picture, err := ChainDataService.GetValidatorAvatar(ctx, rawValidator.Description.Identity, false)

	if err != nil {
		fmt.Println("Get validator avatar failed", err.Error())
	}

	avgAPY, err := service.GetStakeAPY(ctx)
	if err != nil {
		return nil, err
	}

	rate := utils.NewBigFloatFromFloat64(1).Sub(big.NewFloat(commission))

	validator := ValidatorDetail{
		Name:       rawValidator.Description.Moniker,
		Address:    rawValidator.OperatorAddress,
		Stake:      rawValidator.Tokens.Int.String(),
		Voting:     int(votingPower.Int64()),
		Commission: int(commission * 10000),
		Identity:   rawValidator.Description.Identity,
		Picture:    picture,
		APY:        rate.Mul(avgAPY).Mul(big.NewFloat(10000)).Val().String(),
		Active:     active,
	}

	if len(uptimes) == 1 {
		validator.Uptime = uptimes[0].Count * 10000 / blockWindow
	}

	return &validator, nil
}

// func (service *validatortService) SetAccessToken(c *gin.Context, token string) {
// 	prefix := "Bearer "
// 	c.Request.Header.Set("Authorization", prefix+token)
// }
