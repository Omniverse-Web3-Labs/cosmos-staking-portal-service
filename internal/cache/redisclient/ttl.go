package redisclient

import "time"

var (
	RawValidatorsTTL           = 1 * time.Minute
	RawInflationTTL            = 1 * time.Hour
	RawSupplyTTL               = 10 * time.Minute
	RawSlashingParamsTTL       = 1 * time.Hour
	RawStakingParamsTTL        = 1 * time.Hour
	RawDistributionParamsTTL   = 1 * time.Hour
	ValidatorAvatarTTL         = 24 * time.Hour
	RawBlockTTL                = 1 * time.Hour
	RawDelegationsTTL          = 1 * time.Minute
	RawUnbondingDelegationsTTL = 1 * time.Minute
	RawDelegatorRewardsTTL     = 1 * time.Minute
	RawRedelegationsTTL        = 1 * time.Minute
	RawPoolTTL                 = 10 * time.Minute
	RawDelegation              = 5 * time.Minute
	RPCSetNux                  = 5 * time.Second
)
