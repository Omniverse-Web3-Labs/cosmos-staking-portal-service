package router

import (
	"app/internal/controllers/v1/delegator"
	"app/internal/controllers/v1/login"
	"app/internal/controllers/v1/user"
	"app/internal/controllers/v1/validator"

	"github.com/gin-gonic/gin"
)

func RegisterApiV1(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/login/wallet", login.WalletAction)
	routerGroup.POST("/login/cosmos", login.CosmosWalletAction)
	routerGroup.POST("/user/challenge", user.ChallengeAction)
	routerGroup.POST("/login/id", login.IDAction)
	// validator
	routerGroup.POST("/validator/list", validator.ListAction)
	routerGroup.POST("/validator/validators", validator.BatchAction)
	routerGroup.POST("/validator/overview", validator.OverviewAction)

	// validator detail
	routerGroup.POST("/validator/detail", validator.DetailAction)
	routerGroup.POST("/validator/delegation", validator.DelegationAction)
	routerGroup.POST("/validator/commission", validator.CommissionAction)
	routerGroup.POST("/validator/power_event", validator.PowerEventAction)

	// network overview

	// stake/delegator
	routerGroup.POST("/delegator/detail", delegator.DetailAction)
	routerGroup.POST("/delegator/delegation", delegator.UserDelegationAction)
	routerGroup.POST("/delegator/staking_params", delegator.StakingParamsAction)

	// reward
	routerGroup.POST("/delegator/reward_daily", delegator.RewardDailyAction)
	routerGroup.POST("/delegator/reward_rate", delegator.RewardRateAction)
	routerGroup.POST("/delegator/reward_withdrawal", delegator.RewardWithdrawalAction)
	routerGroup.POST("/delegator/reward_total", delegator.RewardTotalAction)
}

func RegisterAuthApiV1(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/user/info", user.InfoAction)
}
