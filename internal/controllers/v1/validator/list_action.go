package validator

import (
	"app/internal/http/resp"
	"app/internal/service"
	"app/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ValidatorListRequest struct {
	Search string `json:"search"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

type ValidatorListResponse struct {
	Total      int                  `json:"total"`
	More       bool                 `json:"more"`
	Offset     int                  `json:"offset"`
	Validators []*service.Validator `json:"validators"`
}

func ListAction(c *gin.Context) {
	var data any
	var request ValidatorListRequest
	err := func() error {
		if err := c.ShouldBindJSON(&request); err != nil {
			return err
		}
		paginater := utils.NewPaginator(request.Page, request.Size)
		desc := false
		if request.Order == "desc" {
			desc = true
		}
		validators, total, err := service.ValidatorService.GetValidatorList(c.Request.Context(), request.Sort, desc, paginater.Offset(), paginater.Limit(), request.Search)
		if err != nil {
			return err
		}
		paginater.SetList(validators)
		paginater.SetTotal(int64(total))
		data = paginater.ToMap()

		// data = ValidatorListResponse{
		// 	Total:      total,
		// 	More:       false,
		// 	Offset:     0,
		// 	Validators: validators,
		// 	Current:    paginater.Current,
		// }
		return nil
	}()
	resp.Finish(c, err, data)
}
