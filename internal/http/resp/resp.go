package resp

import (
	"app/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, utils.JsonCamelCase{Value: data})
}

func Fail(c *gin.Context, err error) {
	code := -1
	httpStatus := http.StatusBadRequest
	if errors.Is(err, ErrInvalidToken) {
		httpStatus = http.StatusUnauthorized
	}
	c.Set("message", err.Error())
	c.Set("code", code)
	c.Error(err)
	c.JSON(httpStatus, gin.H{
		"code":  code,
		"error": err.Error(),
	})
}

func FailWithCode(c *gin.Context, err error, code int) {
	c.Set("message", err.Error())
	c.Set("code", code)
	c.JSON(http.StatusBadRequest, gin.H{
		"code":  code,
		"error": err.Error(),
	})
}

func Finish(c *gin.Context, err error, data any) {
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, data)
}
