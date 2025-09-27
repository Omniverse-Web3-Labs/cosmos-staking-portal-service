package login

import (
	"app/internal/http/resp"

	"github.com/gin-gonic/gin"
)

// 编写获取请求的结构体
type IDRequest struct {
	ID int64 `json:"id"`
}

func IDAction(c *gin.Context) {
	//实例化请求参数
	var request IDRequest
	//返回数据
	var data any
	err := func() error {
		// 解析post的json参数，并且写入request对象
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		data = gin.H{
			"id": request.ID,
		}
		return nil
	}()
	resp.Finish(c, err, data)
}
