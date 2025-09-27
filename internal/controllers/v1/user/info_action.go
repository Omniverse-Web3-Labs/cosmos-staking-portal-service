package user

import (
	"app/internal/db/postgres/query"
	"app/internal/http/resp"
	"app/internal/service"
	"app/pkg/utils/helper"

	"github.com/gin-gonic/gin"
)

func InfoAction(c *gin.Context) {
	ctx := c.Request.Context()
	var data any
	err := func() error {
		id := service.AuthService.CurrentAuthID(c)
		data = gin.H{
			"id": id,
		}
		user, err := query.User.WithContext(ctx).Where(query.User.ID.Eq(helper.Int64Val(id))).First()
		if err != nil {
			return err
		}
		data = user
		return nil
	}()
	resp.Finish(c, err, data)
}
