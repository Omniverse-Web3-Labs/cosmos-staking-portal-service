package redisclient

var (
	KeyRawValidators            = "validators"
	KeyValidatorsLock           = "lock:validator_list"
	KeyRawSupply                = "supply"
	KeyRawPool                  = "pool"
	KeyRawSlashingParams        = "slashing_params"
	KeyRawDistributionParams    = "distribution_params"
	KeyRawStakingParams         = "staking_params"
	KeyRawBlock                 = "block:%d"                 // block:<height>
	KeyRawDelegations           = "delegations:%s"           // delegations:<delegator>
	KeyRawUnbondingDelegations  = "unbonding_delegations:%s" // unbonding_delegations:<delegator>
	KeyValidatorAvatar          = "validator_avatar:%s"      // validator_avatar:<identity>
	KeyRawDelegatorRewards      = "delegator_rewards:%s"     // delegator_rewards:<delegator>
	KeyRawRedelegations         = "redelegations:%s"         // redelegations:<delegator>
	KeyRawInflation             = "inflation"
	KeySupplyLock               = "lock:supply"
	KeyInflationLock            = "lock:inflation"
	KeyPoolLock                 = "lock:pool"
	KeySlashingParamsLock       = "lock:slashing_params"
	KeyDistributionParamsLock   = "lock:distribution_params"
	KeyStakingParamsLock        = "lock:staking_params"
	KeyBlockLock                = "lock:block:%d"
	KeyDelegationsLock          = "lock:delegations:%s"
	KeyUnbondingDelegationsLock = "lock:unbonding_delegations:%s"
	KeyValidatorAvatarLock      = "lock:validator_avatar:%s"
	KeyRawDelegatorRewardsLock  = "lock:delegator_rewards:%s"
	KeyRawRedelegationsLock     = "lock:redelegations:%s"
	KeyRawDelegators            = "delegation:%s:%s:%d:%d:%t"      // delegation:<validator>:<sort_key>:<offset>:<limit>:<desc>
	KeyDelegationLock           = "lock:delegation:%s:%s:%d:%d:%t" // lock:delegation:<validator>:<sort_key>:<offset>:<limit>:<desc>
)
