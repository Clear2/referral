package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/logger"
	"github.com/keep/sunny/pkg/response"
	"github.com/keep/sunny/pkg/utils"
)

type Controller struct {
	service    *Service
	repository *Repository
	logger     logger.Logger
}

func NewController(service *Service, repository *Repository, log logger.Logger) *Controller {
	return &Controller{service: service, repository: repository, logger: log}
}

func bindID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		_ = c.Error(appErrors.ErrBadRequest)
		return 0, false
	}
	return id, true
}
func bindJSON[T any](c *gin.Context) (T, bool) {
	var input T
	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return input, false
	}
	return input, true
}
func (controller *Controller) audit(c *gin.Context, id int, action string) error {
	operator, err := utils.MustGetUser(c)
	if err != nil {
		return fmt.Errorf("resolve audit operator: %w", err)
	}
	return controller.repository.Audit(c.Request.Context(), operator, action, id, requestid.Get(c), c.ClientIP())
}
func (controller *Controller) changed(c *gin.Context, id int, action string, err error) {
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err = controller.audit(c, id, action); err != nil {
		controller.logger.Errorw("write user audit log failed", "action", action, "target_id", id, "error", err)
	}
	response.WriteSuccess(c, gin.H{"updated": true})
}

func (controller *Controller) List(c *gin.Context) {
	var input ListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return
	}
	out, err := controller.service.List(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}
func (controller *Controller) Create(c *gin.Context) {
	input, ok := bindJSON[CreateInput](c)
	if !ok {
		return
	}
	out, err := controller.service.Create(c.Request.Context(), input)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err = controller.audit(c, out.ID, "CREATE"); err != nil {
		controller.logger.Errorw("write user audit log failed", "action", "CREATE", "target_id", out.ID, "error", err)
	}
	response.WriteSuccess(c, out)
}
func (controller *Controller) Get(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	out, err := controller.service.Get(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}
func (controller *Controller) Update(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	input, ok := bindJSON[UpdateInput](c)
	if ok {
		controller.changed(c, id, "UPDATE", controller.service.Update(c.Request.Context(), id, input))
	}
}
func (controller *Controller) SetStatus(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	input, ok := bindJSON[StatusInput](c)
	if ok {
		operatorID, err := utils.MustGetUser(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		controller.changed(c, id, "STATUS", controller.service.SetStatus(c.Request.Context(), operatorID, id, input.Enabled))
	}
}
func (controller *Controller) SetRoles(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	input, ok := bindJSON[RolesInput](c)
	if ok {
		operatorID, err := utils.MustGetUser(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		controller.changed(c, id, "ASSIGN", controller.service.SetRoles(c.Request.Context(), operatorID, id, input.RoleIDs))
	}
}
func (controller *Controller) ResetPassword(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	input, ok := bindJSON[ResetPasswordInput](c)
	if ok {
		controller.changed(c, id, "RESET_PASSWORD", controller.service.ResetPassword(c.Request.Context(), id, input))
	}
}
