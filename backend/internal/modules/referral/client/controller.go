package client

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/internal/modules/auth"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/response"
)

type sessionStarter interface {
	StartSession(*gin.Context, int) error
}

type Controller struct {
	service  Service
	sessions sessionStarter
}

func NewController(service Service, sessions *auth.Manager) *Controller {
	return &Controller{service: service, sessions: sessions}
}

type userIDParams struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// GenerateInvitation godoc
// @Summary Generate or return a user's invitation link
// @Tags Referrals
// @Produce json
// @Success 200 {object} InvitationAPIResponse
// @Router /api/v1/users/{id}/referral-code [post]
func (ctl *Controller) GenerateInvitation(c *gin.Context) {
	var params userIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return
	}
	code, err := ctl.service.GenerateInvitation(c.Request.Context(), params.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	response.WriteSuccess(c, InvitationView{Code: code, URL: fmt.Sprintf("%s://%s/ref/%s", scheme, c.Request.Host, code)})
}

// Register godoc
// @Summary Register a new user through an invitation
// @Tags Referrals
// @Accept json
// @Produce json
// @Success 200 {object} RegistrationAPIResponse
// @Router /api/v1/referrals/register [post]
func (ctl *Controller) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return
	}
	out, err := ctl.service.Register(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err = ctl.sessions.StartSession(c, out.Invitee.ID); err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}

// Dashboard godoc
// @Summary Get referral statistics, history and credit ledger
// @Tags Referrals
// @Produce json
// @Success 200 {object} DashboardAPIResponse
// @Router /api/v1/users/{id}/referral-dashboard [get]
func (ctl *Controller) Dashboard(c *gin.Context) {
	var params userIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return
	}
	out, err := ctl.service.Dashboard(c.Request.Context(), params.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}
