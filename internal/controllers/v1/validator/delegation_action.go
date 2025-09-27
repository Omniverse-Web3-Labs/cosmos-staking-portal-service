package validator

import (
	"app/internal/http/resp"
	"app/internal/service"
	"app/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DelegationRequest struct {
	Address string `json:"address"`
	Sort    string `json:"sort"`
	Order   string `json:"order"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type DelegationResponse struct {
	Total      int                 `json:"total"`
	More       bool                `json:"more"`
	Offset     int                 `json:"offset"`
	Delegators []service.Delegator `json:"delegators"`
}

func DelegationAction(c *gin.Context) {
	var data any
	var request DelegationRequest
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}

		desc := false
		if request.Order == "desc" {
			desc = true
		}

		paginator := utils.NewPaginator(request.Page, request.Size)
		delegators, total, err := service.ValidatorService.GetDelegators(c.Request.Context(), request.Address, request.Sort, paginator.Offset(), paginator.Limit(), desc)
		if err != nil {
			return err
		}

		more := false
		offset := paginator.Offset() + len(delegators)
		if offset < total {
			more = true
		}
		paginator.SetTotal(int64(total))
		paginator.HasMore = more
		paginator.List = delegators
		data = paginator.ToMap()
		/*
			data = DelegationResponse{
				Total:      total,
				More:       more,
				Offset:     offset,
				Delegators: delegators,
			}
		*/
		return nil
	}()
	resp.Finish(c, err, data)
}
