package service

import (
	"app/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

var RequestService requestService

type requestService struct {
}

func (requestService) ShouldBindJSON(c *gin.Context, obj any) error {
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	return nil
}

func (service *requestService) GetAccessToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	prefix := "Bearer "
	if utils.StringUtil.StartsWith(token, prefix) {
		return strings.TrimPrefix(token, prefix)
	}
	return ""
}

func (service *requestService) SetAccessToken(c *gin.Context, token string) {
	prefix := "Bearer "
	c.Request.Header.Set("Authorization", prefix+token)
}
