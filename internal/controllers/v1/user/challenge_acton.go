package user

import (
	"app/internal/http/resp"
	"app/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChallengeReqeust struct {
	Address string `json:"address" binding:"required"`
}

func ChallengeAction(c *gin.Context) {
	ctx := c.Request.Context()
	//实例化请求参数
	var request ChallengeReqeust
	//返回数据
	var data any
	err := func() error {
		// 解析post的json参数，并且写入request对象
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		challenge, err := service.UserService.GetChallenge(ctx, strings.ToLower(request.Address))
		if err != nil {
			return err
		}

		data = gin.H{
			"challenge": challenge,
		}
		return nil
	}()
	resp.Finish(c, err, data)
}
