package service

import (
	"app/internal/cache/redisclient"
	"app/internal/config"
	"app/internal/keybase"
	"app/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ChainDataService chainDataService

type BigIntFromString struct {
	*big.Int
}

func (b *BigIntFromString) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"") // 去掉引号
	bInt := new(big.Int)
	bInt, success := bInt.SetString(str, 10)
	if !success {
		return fmt.Errorf("invalid big integer string: %s", str)
	}
	b.Int = bInt
	return nil
}

func (b BigIntFromString) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

type BigFloatFromString struct {
	*big.Float
}

func (b *BigFloatFromString) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	bigFloat := new(big.Float)

	_, success := bigFloat.SetString(str)
	if !success {
		return fmt.Errorf("invalid big float string: %s", str)
	}

	b.Float = bigFloat
	return nil
}

type Float64FromString float64

func (f *Float64FromString) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1] // 去掉引号
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*f = Float64FromString(val)
	return nil
}

func (f Float64FromString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatFloat(float64(f), 'f', -1, 64))
}

type Int64FromString int64

func (f *Int64FromString) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1] // 去掉引号
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*f = Int64FromString(val)
	return nil
}

func (i Int64FromString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

type RawDelegation struct {
	Delegation []RawDelegator `json:"delegation_responses"`
	Pagination RawPagination  `json:"pagination"`
}

type RawUnbondingDelegator struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Entries          []struct {
		CreationHeight          Int64FromString  `json:"creation_height"`
		CompletionTime          time.Time        `json:"completion_time"`
		InitialBalance          BigIntFromString `json:"initial_balance"`
		Balance                 BigIntFromString `json:"balance"`
		UnbondingID             Int64FromString  `json:"unbonding_id"`
		UnbondingOnHoldRefCount Int64FromString  `json:"unbonding_on_hold_ref_count"`
	} `json:"entries"`
}

type RawDelegatorRewards struct {
	Rewards []struct {
		ValidatorAddress string  `json:"validator_address"`
		Reward           []Token `json:"reward"`
	} `json:"rewards"`
	Total []Token `json:"total"`
}

type Token struct {
	Denom  string           `json:"denom"`
	Amount BigIntFromString `json:"amount"`
}

type RawRedelegations struct {
	Responses []struct {
		Redelegation struct {
			Delegator    string `json:"delegator_address"`
			ValidatorSrc string `json:"validator_src_address"`
			ValidatorDst string `json:"validator_dst_address"`
		} `json:"redelegation"`
		Entries []struct {
			RedelegationEntry struct {
				CreationHeight int64              `json:"creation_height"`
				CompletionTime time.Time          `json:"completion_time"`
				InitialBalance BigIntFromString   `json:"initial_balance"`
				SharesDst      BigFloatFromString `json:"shares_dst"`
				UnbondingID    int64              `json:"unbonding_id"`
			} `json:"redelegation_entry"`
		} `json:"entries"`
	} `json:"redelegation_responses"`
	Pagination RawPagination `json:"pagination"`
}

type RawUnbondingDelegation struct {
	Unbondings []RawUnbondingDelegator `json:"unbonding_responses"`
	Pagination RawPagination           `json:"pagination"`
}

type RawDelegator struct {
	Delegation struct {
		DelegatorAddress string             `json:"delegator_address"`
		ValidatorAddress string             `json:"validator_address"`
		Shares           BigFloatFromString `json:"shares"`
	} `json:"delegation"`
	Balance struct {
		Denom  string           `json:"denom"`
		Amount BigIntFromString `json:"amount"`
	} `json:"balance"`
}

type RawValidators struct {
	Validators []RawValidator `json:"validators"`
	Pagination RawPagination  `json:"pagination"`
}

type RawSingleValidator struct {
	SingleValidator RawValidator `json:"validator"`
}

type RawPagination struct {
	NextKey string          `json:"next_key"`
	Total   Int64FromString `json:"total"`
}

type RawPool struct {
	PoolData struct {
		NotBondedTokens BigIntFromString `json:"not_bonded_tokens"`
		BondedTokens    BigIntFromString `json:"bonded_tokens"`
	} `json:"pool"`
}

type RawSupply struct {
	Supply []RawSingleSupply `json:"supply"`
}

type RawSingleSupply struct {
	Denom  string           `json:"denom"`
	Amount BigIntFromString `json:"amount"`
}

type RawInflation struct {
	Inflation Float64FromString `json:"inflation"`
}

type RawValidator struct {
	OperatorAddress string             `json:"operator_address"`
	Jailed          bool               `json:"jailed"`
	Status          string             `json:"status"`
	Tokens          BigIntFromString   `json:"tokens"`
	DelegatorShares BigFloatFromString `json:"delegator_shares"`
	Description     struct {
		Moniker          string `json:"moniker"`
		Identity         string `json:"identity"`
		Website          string `json:"website"`
		SecurityContract string `json:"security_contact"`
		Details          string `json:"details"`
	} `json:"description"`
	Commission struct {
		CommissionRates struct {
			Rate          Float64FromString `json:"rate"`
			MaxRate       Float64FromString `json:"max_rate"`
			MaxChangeRate Float64FromString `json:"max_change_rate"`
		} `json:"commission_rates"`
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	UnbondingHeight   Int64FromString `json:"unbonding_height"`
	UnbondingTime     time.Time       `json:"unbonding_time"`
	MinSelfDelegation Int64FromString `json:"min_self_delegation"`
}

type RawBlock struct {
	Block struct {
		Header struct {
			Height Int64FromString `json:"height"`
			Time   time.Time       `json:"time"`
		} `json:"header"`
	} `json:"block"`
}

type RawStakingParams struct {
	Params struct {
		UnbondingTime string `json:"unbonding_time"`
		MaxEntries    int64  `json:"max_entries"`
	} `json:"params"`
}

type SlashingParams struct {
	Params struct {
		SignedBlocksWindow      Int64FromString   `json:"signed_blocks_window"`
		MinSignedPerWindow      Float64FromString `json:"min_signed_per_window"`
		DowntimeJailDuration    string            `json:"downtime_jail_duration"`
		SlashFractionDoubleSign Float64FromString `json:"slash_fraction_double_sign"`
		SlashFractionDowntime   Float64FromString `json:"slash_fraction_downtime"`
	} `json:"params"`
}

type DistributionParams struct {
	Params struct {
		CommunityTax        Float64FromString `json:"community_tax"`
		BaseProposerReward  Float64FromString `json:"base_proposer_reward"`
		BonusProposerReward Float64FromString `json:"bonus_proposer_reward"`
		WithdrawAddrEnabled bool              `json:"withdraw_addr_enabled"`
	} `json:"params"`
}

type chainDataService struct {
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/validators
func (service *chainDataService) GetValidators(ctx context.Context) ([]RawValidator, error) {
	data, err := redisclient.RedisClient.HGet(ctx, redisclient.KeyRawValidators, "raw_list").Bytes()
	if err != nil {
		fmt.Printf("redis get %s error %s\n", redisclient.KeyRawValidators, err.Error())
	}

	var validators RawValidators

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeyValidatorsLock, 1, redisclient.RPCSetNux).Result()

		if err != nil {
			fmt.Println("Redis set nx validator list error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/validators", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(data, &validators)
			if err != nil {
				return nil, err
			}

			pipe := redisclient.RedisClient.TxPipeline()
			pipe.HSet(ctx, redisclient.KeyRawValidators, "raw_list", data)
			for _, v := range validators.Validators {
				vData, err := json.Marshal(v)
				if err != nil {
					fmt.Println("validator data marshal failed", err.Error())
					continue
				}
				pipe.HSet(ctx, redisclient.KeyRawValidators, v.OperatorAddress, vData)
			}
			pipe.Expire(ctx, redisclient.KeyRawValidators, redisclient.RawValidatorsTTL)
			_, err = pipe.Exec(ctx)
			if err != nil {
				fmt.Println("Redis Get Validators pipe exec failed", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeyValidatorsLock)
		} else {
			checkInterval := 100 * time.Millisecond
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.HGet(ctx, redisclient.KeyRawValidators, "raw_list").Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						err = json.Unmarshal(data, &validators)
						if err != nil {
							fmt.Println("tick got data unmarshal err", err.Error())
							return nil, err
						}
						fmt.Println("validator list tick got data", validators)
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw validators timeout")
				}
			}
		}
	} else {
		err = json.Unmarshal(data, &validators)
		if err != nil {
			return nil, err
		}
	}

	return validators.Validators, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/validator/ethmvaloper1clnaafa38d785mstfw9c29hxwg0tdjmjmuw07v
func (service *chainDataService) GetValidator(ctx context.Context, validatorAddress string) (*RawValidator, error) {
	data, err := redisclient.RedisClient.HGet(ctx, redisclient.KeyRawValidators, validatorAddress).Bytes()
	if err != nil {
		fmt.Printf("redis get %s:%s error %s\n", redisclient.KeyRawValidators, validatorAddress, err.Error())
	}

	var validator RawValidator
	if data == nil {
		_, err = service.GetValidators(ctx)
		if err != nil {
			return nil, err
		}

		data, err = redisclient.RedisClient.HGet(ctx, redisclient.KeyRawValidators, validatorAddress).Bytes()
		if err != nil {
			fmt.Printf("redis get %s:%s error %s\n", redisclient.KeyRawValidators, validatorAddress, err.Error())
			return nil, err
		} else {
			if data == nil {
				fmt.Println("get validator error, no data")
				return nil, fmt.Errorf("get validator error, no data")
			}
		}
	}

	err = json.Unmarshal(data, &validator)
	if err != nil {
		return nil, err
	}

	return &validator, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/validators/ethmvaloper1clnaafa38d785mstfw9c29hxwg0tdjmjmuw07v/delegations
func (service *chainDataService) GetDelegators(ctx context.Context, validator string, key string, offset int, limit int, desc bool) ([]RawDelegator, error) {
	redisKey := fmt.Sprintf(redisclient.KeyRawDelegators, validator, key, offset, limit, desc)
	data, err := redisclient.RedisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		fmt.Printf("redis get %s error %s\n", redisKey, err.Error())
	}

	if data == nil {
		// get lock
		redisLockKey := fmt.Sprintf(redisclient.KeyDelegationLock, validator, key, offset, limit, desc)
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisLockKey, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx pool error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/validators/%s/delegations?key=%s&offset=%d&limit=%d&reverse=%t",
				config.GetConfig().Endpoint.Rest, validator, key, offset, limit, desc)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err := redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.RawDelegation).Err()
			if err != nil {
				fmt.Println("Redis set pool error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisLockKey)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					redisKey := fmt.Sprintf(redisclient.KeyRawDelegators, validator, key, offset, limit, desc)
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw delegations timeout")
				}
			}
		}
	}

	var delegation RawDelegation
	err = json.Unmarshal(data, &delegation)
	if err != nil {
		return nil, err
	}

	return delegation.Delegation, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/delegators/ethmvaloper1clnaafa38d785mstfw9c29hxwg0tdjmjmuw07v/unbonding_delegations
func (service *chainDataService) GetUnbondingDelegations(ctx context.Context, delegator string) ([]RawUnbondingDelegator, error) {
	redisKey := fmt.Sprintf(redisclient.KeyRawUnbondingDelegations, delegator)
	data, err := redisclient.RedisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		fmt.Printf("GetUnbondingDelegations %s error %s\n", delegator, err.Error())
	}

	if data == nil {
		// get lock
		redisLockKey := fmt.Sprintf(redisclient.KeyUnbondingDelegationsLock, delegator)
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisLockKey, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Printf("Redis set nx unbonding delegations %s error %s\n", delegator, err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/delegators/%s/unbonding_delegations", config.GetConfig().Endpoint.Rest, delegator)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.RawUnbondingDelegationsTTL).Err()
			if err != nil {
				fmt.Printf("Redis set unbonding delegations %s error %s\n", delegator, err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisLockKey)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw unbonding delegations %s timeout", delegator)
				}
			}
		}
	}

	var unbondingDelegation RawUnbondingDelegation
	err = json.Unmarshal(data, &unbondingDelegation)
	if err != nil {
		return nil, err
	}

	return unbondingDelegation.Unbondings, nil
}

// real-time data
// curl http://localhost:1317/cosmos/distribution/v1beta1/delegators/omni195tf7tuvxagrsq7seh8579qrjvxc6j9wgnpz86/rewards
func (service *chainDataService) GetDelegatorRewards(ctx context.Context, delegator string) (*RawDelegatorRewards, error) {
	redisKey := fmt.Sprintf(redisclient.KeyRawDelegatorRewards, delegator)
	data, err := redisclient.RedisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		fmt.Printf("GetDelegatorRewards %s error %s\n", delegator, err.Error())
	}

	if data == nil {
		// get lock
		redisLockKey := fmt.Sprintf(redisclient.KeyRawDelegatorRewardsLock, delegator)
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisLockKey, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Printf("Redis set nx delegator rewards %s error %s\n", delegator, err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/distribution/v1beta1/delegators/%s/rewards", config.GetConfig().Endpoint.Rest, delegator)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.RawDelegatorRewardsTTL).Err()
			if err != nil {
				fmt.Printf("Redis set delegator rewards %s error %s\n", delegator, err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisLockKey)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw delegator rewards %s timeout", delegator)
				}
			}
		}
	}

	var delegatorRewards RawDelegatorRewards
	err = json.Unmarshal(data, &delegatorRewards)
	if err != nil {
		return nil, err
	}

	return &delegatorRewards, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/delegators/omni195tf7tuvxagrsq7seh8579qrjvxc6j9wgnpz86/redelegations
func (service *chainDataService) GetRedelegations(ctx context.Context, delegator string) (*RawRedelegations, error) {
	redisKey := fmt.Sprintf(redisclient.KeyRawRedelegations, delegator)
	data, err := redisclient.RedisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		fmt.Printf("GetRedelegations %s error %s\n", delegator, err.Error())
	}

	if data == nil {
		// get lock
		redisLockKey := fmt.Sprintf(redisclient.KeyRawRedelegationsLock, delegator)
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisLockKey, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Printf("Redis set nx redelegations %s error %s\n", delegator, err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/delegators/%s/redelegations", config.GetConfig().Endpoint.Rest, delegator)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.RawRedelegationsTTL).Err()
			if err != nil {
				fmt.Printf("Redis set redelegations %s error %s\n", delegator, err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisLockKey)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw delegator rewards %s timeout", delegator)
				}
			}
		}
	}

	var redelegations RawRedelegations
	err = json.Unmarshal(data, &redelegations)
	if err != nil {
		return nil, err
	}

	return &redelegations, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/delegations/ethmvaloper1clnaafa38d785mstfw9c29hxwg0tdjmjmuw07v
func (service *chainDataService) GetDelegations(ctx context.Context, delegator string) ([]RawDelegator, error) {
	redisKey := fmt.Sprintf(redisclient.KeyRawDelegations, delegator)
	data, err := redisclient.RedisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		fmt.Printf("GetDelegations %s error %s\n", delegator, err.Error())
	}

	if data == nil {
		// get lock
		redisLockKey := fmt.Sprintf(redisclient.KeyDelegationsLock, delegator)
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisLockKey, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Printf("Redis set nx delegations %s error %s\n", delegator, err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/delegations/%s", config.GetConfig().Endpoint.Rest, delegator)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.RawDelegationsTTL).Err()
			if err != nil {
				fmt.Printf("Redis set delegations %s error %s\n", delegator, err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisLockKey)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw delegations %s timeout", delegator)
				}
			}
		}
	}

	var delegation RawDelegation
	err = json.Unmarshal(data, &delegation)
	if err != nil {
		return nil, err
	}

	return delegation.Delegation, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/pool
func (service *chainDataService) GetPool(ctx context.Context, update bool) (*RawPool, error) {
	var data []byte
	var err error
	if !update {
		data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawPool).Bytes()
		if err != nil {
			fmt.Println("GetPool error", err.Error())
		}
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeyPoolLock, 1, redisclient.RPCSetNux).Result()

		if err != nil {
			fmt.Println("Redis set nx pool error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/pool", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err := redisclient.RedisClient.Set(ctx, redisclient.KeyRawPool, data, redisclient.RawPoolTTL).Err()
			if err != nil {
				fmt.Println("Redis set pool error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeyPoolLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawPool).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw pool timeout")
				}
			}
		}
	}

	var pool RawPool
	err = json.Unmarshal(data, &pool)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// real-time data
// curl http://localhost:1317/cosmos/bank/v1beta1/supply
func (service *chainDataService) GetTotalSupply(ctx context.Context, update bool) (*RawSingleSupply, error) {
	var data []byte
	var err error
	if !update {
		data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawSupply).Bytes()
		if err != nil {
			fmt.Println("GetTotalSupply error", err.Error())
		}
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeySupplyLock, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx supply error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/bank/v1beta1/supply", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err := redisclient.RedisClient.Set(ctx, redisclient.KeyRawSupply, data, redisclient.RawSupplyTTL).Err()
			if err != nil {
				fmt.Println("Redis set supply error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeySupplyLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawSupply).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw total supply timeout")
				}
			}
		}
	}

	var supply RawSupply
	err = json.Unmarshal(data, &supply)
	if err != nil {
		return nil, err
	}

	return &supply.Supply[0], nil
}

// real-time data
// curl http://localhost:1317/cosmos/mint/v1beta1/inflation
func (service *chainDataService) GetInflation(ctx context.Context, update bool) (float64, error) {
	var data []byte
	var err error
	if !update {
		data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawInflation).Bytes()
		if err != nil {
			fmt.Println("GetInflation error", err.Error())
		}
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeyInflationLock, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx inflation error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/mint/v1beta1/inflation", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return 0, err
			}

			err := redisclient.RedisClient.Set(ctx, redisclient.KeyRawInflation, data, redisclient.RawInflationTTL).Err()
			if err != nil {
				fmt.Println("Redis set inflation error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeyInflationLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawInflation).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return 0, fmt.Errorf("get raw total supply timeout")
				}
			}
		}
	}

	var inflation RawInflation
	err = json.Unmarshal(data, &inflation)
	if err != nil {
		return 0, err
	}

	return float64(inflation.Inflation), nil
}

// real-time data
// curl http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/{height}
func (service *chainDataService) GetBlock(ctx context.Context, height int64) (*RawBlock, error) {
	redisKey := fmt.Sprintf(redisclient.KeyRawBlock, height)
	data, err := redisclient.RedisClient.Get(ctx, redisKey).Bytes()
	if err != nil {
		fmt.Printf("GetBlock %d error %s\n", height, err.Error())
	}

	if data == nil {
		// get lock
		redisLockKey := fmt.Sprintf(redisclient.KeyBlockLock, height)
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisLockKey, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Printf("Redis set nx block %d error %s\n", height, err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/base/tendermint/v1beta1/blocks/%d", config.GetConfig().Endpoint.Rest, height)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.RawBlockTTL).Err()
			if err != nil {
				fmt.Printf("Redis set block %d error %s\n", height, err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisLockKey)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw slashing params timeout")
				}
			}
		}
	}

	var block RawBlock
	err = json.Unmarshal(data, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// real-time data
// curl http://localhost:1317/cosmos/staking/v1beta1/params
func (service *chainDataService) GetStakingParams(ctx context.Context) (*RawStakingParams, error) {
	data, err := redisclient.RedisClient.Get(ctx, redisclient.KeyRawStakingParams).Bytes()
	if err != nil {
		fmt.Println("GetStakingParams error", err.Error())
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeyStakingParamsLock, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx staking params error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/staking/v1beta1/params", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err = redisclient.RedisClient.Set(ctx, redisclient.KeyRawStakingParams, data, redisclient.RawStakingParamsTTL).Err()
			if err != nil {
				fmt.Println("Redis set staking params error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeyStakingParamsLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawStakingParams).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw slashing params timeout")
				}
			}
		}
	}

	var params RawStakingParams
	err = json.Unmarshal(data, &params)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

// real-time data
// curl http://localhost:1317/cosmos/slashing/v1beta1/params
func (service *chainDataService) GetSlashingParams(ctx context.Context) (*SlashingParams, error) {
	data, err := redisclient.RedisClient.Get(ctx, redisclient.KeyRawSlashingParams).Bytes()
	if err != nil {
		fmt.Println("GetSlashingParams error", err.Error())
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeySlashingParamsLock, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx slashing params error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/slashing/v1beta1/params", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err := redisclient.RedisClient.Set(ctx, redisclient.KeyRawSlashingParams, data, redisclient.RawSlashingParamsTTL).Err()
			if err != nil {
				fmt.Println("Redis set slashing params error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeySlashingParamsLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawSlashingParams).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw slashing params timeout")
				}
			}
		}
	}

	var params SlashingParams
	err = json.Unmarshal(data, &params)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

// real-time data
// curl http://localhost:1317/cosmos/distribution/v1beta1/params
func (service *chainDataService) GetDistributionParams(ctx context.Context) (*DistributionParams, error) {
	data, err := redisclient.RedisClient.Get(ctx, redisclient.KeyRawDistributionParams).Bytes()
	if err != nil {
		fmt.Println("GetDistributionParams error", err.Error())
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeyDistributionParamsLock, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx distribution params error", err.Error())
		} else if setnx {
			url := fmt.Sprintf("%scosmos/distribution/v1beta1/params", config.GetConfig().Endpoint.Rest)

			data, err = utils.HttpUtil.SendRequestWithContext(ctx, url, "GET", nil, nil)

			if err != nil {
				return nil, err
			}

			err := redisclient.RedisClient.Set(ctx, redisclient.KeyRawDistributionParams, data, redisclient.RawDistributionParamsTTL).Err()
			if err != nil {
				fmt.Println("Redis set distribution params error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeyDistributionParamsLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisclient.KeyRawDistributionParams).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return nil, fmt.Errorf("get raw slashing params timeout")
				}
			}
		}
	}

	var params DistributionParams
	err = json.Unmarshal(data, &params)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

func (service *chainDataService) GetValidatorAvatar(ctx context.Context, identity string, update bool) (string, error) {
	matched, err := regexp.MatchString(`^[a-fA-F0-9]{16}$`, identity)
	if err != nil {
		return "", err
	}

	if !matched {
		return "", fmt.Errorf("identity %s format err", identity)
	}

	redisKey := fmt.Sprintf(redisclient.KeyValidatorAvatar, identity)
	var data []byte
	if !update {
		data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
		if err != nil {
			fmt.Printf("GetValidatorAvatar %s error %s\n", identity, err.Error())
		}
	}

	if data == nil {
		// get lock
		setnx, err := redisclient.RedisClient.SetNX(ctx, redisclient.KeyValidatorAvatarLock, 1, redisclient.RPCSetNux).Result()
		if err != nil {
			fmt.Println("Redis set nx validator avatar error", err.Error())
		} else if setnx {
			data, err := keybase.Keybase.GetAvatar(ctx, identity)
			if err != nil {
				return "", err
			}

			err = redisclient.RedisClient.Set(ctx, redisKey, data, redisclient.ValidatorAvatarTTL).Err()
			if err != nil {
				fmt.Println("Redis set validator avatar error", err.Error())
			}
			redisclient.RedisClient.Del(ctx, redisclient.KeyValidatorAvatarLock)
		} else {
			checkInterval := 100 * time.Millisecond // 0.1秒
			totalTime := redisclient.RPCSetNux

			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			timeout := time.After(totalTime)

		outer:
			for {
				select {
				case <-ticker.C:
					data, err = redisclient.RedisClient.Get(ctx, redisKey).Bytes()
					if err != nil {
						continue
					}

					if data != nil {
						break outer
					}
				case <-timeout:
					return "", fmt.Errorf("get validator avatar timeout")
				}
			}
		}
	}

	return string(data), nil
}
