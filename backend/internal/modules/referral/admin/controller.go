package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/response"
)

type Controller struct{ service *Service }

func NewController(service *Service) *Controller { return &Controller{service: service} }

func bindList(c *gin.Context) (ListInput, bool) {
	var input ListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return input, false
	}
	return input, true
}

// Referrals godoc
// @Summary List all referral rewards
// @Tags Admin Referrals
// @Produce json
// @Success 200 {object} ReferralListAPIResponse
// @Router /api/v1/admin/referrals [get]
func (ctl *Controller) Referrals(c *gin.Context) {
	input, ok := bindList(c)
	if !ok {
		return
	}
	out, err := ctl.service.Referrals(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}

// CreditTransactions godoc
// @Summary List all referral credit transactions
// @Tags Admin Referrals
// @Produce json
// @Success 200 {object} CreditTransactionListAPIResponse
// @Router /api/v1/admin/credit-transactions [get]
func (ctl *Controller) CreditTransactions(c *gin.Context) {
	input, ok := bindList(c)
	if !ok {
		return
	}
	out, err := ctl.service.CreditTransactions(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}

// Stats godoc
// @Summary Get global referral reward statistics
// @Tags Admin Referrals
// @Produce json
// @Success 200 {object} StatsAPIResponse
// @Router /api/v1/admin/referral-stats [get]
func (ctl *Controller) Stats(c *gin.Context) {
	out, err := ctl.service.Stats(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}
