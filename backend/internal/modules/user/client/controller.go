package client

import (
	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/pkg/response"
	"github.com/keep/sunny/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

// GetMe godoc
// @Summary     Get current user
// @Description Retrieve the authenticated user's information
// @ID          getCurrentUser
// @Tags        Users
// @Produce     json
// @Success     200 {object} response.UserAPIResponse "Successfully retrieved user"
// @Failure     401 {object} errors.APIError "Unauthorized"
// @Failure     404 {object} errors.APIError "User not found"
// @Failure     500 {object} errors.APIError "Internal server error"
// @Router      /api/v1/users/me [get]
func (u *UserController) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	id, err := utils.MustGetUser(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	span.SetAttributes(attribute.Bool("has_user_id", true))
	user, err := u.userService.GetByID(ctx, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.WriteSuccess(c, user)
}
