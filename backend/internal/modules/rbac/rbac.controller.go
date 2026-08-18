package rbac

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/response"
	"github.com/keep/sunny/pkg/utils"
)

type Controller struct {
	service *Service
	engine  *gin.Engine
}

func NewController(service *Service, engine *gin.Engine) *Controller {
	return &Controller{service: service, engine: engine}
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
	var in T
	if err := c.ShouldBindJSON(&in); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, appErrors.ValidationMessage(err)))
		return in, false
	}
	return in, true
}
func done(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, gin.H{"updated": true})
}
func (ctl *Controller) audited(c *gin.Context, err error, action, target string, id int) {
	if err != nil {
		done(c, err)
		return
	}
	operatorID, authErr := utils.MustGetUser(c)
	if authErr != nil {
		done(c, authErr)
		return
	}
	if auditErr := ctl.service.repo.Audit(c.Request.Context(), operatorID, action, target, id, requestid.Get(c), c.ClientIP()); auditErr != nil {
		done(c, publicError(auditErr))
		return
	}
	done(c, nil)
}

func (ctl *Controller) Snapshot(c *gin.Context) {
	out, err := ctl.service.Snapshot(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}
func (ctl *Controller) CreateRole(c *gin.Context) {
	in, ok := bindJSON[RoleInput](c)
	if ok {
		ctl.audited(c, ctl.service.CreateRole(c.Request.Context(), in), "CREATE", "ROLE", 0)
	}
}
func (ctl *Controller) UpdateRole(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[RoleInput](c)
	if ok {
		ctl.audited(c, ctl.service.UpdateRole(c.Request.Context(), id, in), "UPDATE", "ROLE", id)
	}
}
func (ctl *Controller) DeleteRole(c *gin.Context) {
	id, ok := bindID(c)
	if ok {
		ctl.audited(c, ctl.service.DeleteRole(c.Request.Context(), id), "DELETE", "ROLE", id)
	}
}
func (ctl *Controller) CreatePermission(c *gin.Context) {
	in, ok := bindJSON[PermissionInput](c)
	if ok {
		ctl.audited(c, ctl.service.CreatePermission(c.Request.Context(), in), "CREATE", "PERMISSION", 0)
	}
}
func (ctl *Controller) UpdatePermission(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[PermissionInput](c)
	if ok {
		ctl.audited(c, ctl.service.UpdatePermission(c.Request.Context(), id, in), "UPDATE", "PERMISSION", id)
	}
}
func (ctl *Controller) DeletePermission(c *gin.Context) {
	id, ok := bindID(c)
	if ok {
		ctl.audited(c, ctl.service.DeletePermission(c.Request.Context(), id), "DELETE", "PERMISSION", id)
	}
}
func (ctl *Controller) CreateMenu(c *gin.Context) {
	in, ok := bindJSON[MenuInput](c)
	if ok {
		ctl.audited(c, ctl.service.CreateMenu(c.Request.Context(), in), "CREATE", "MENU", 0)
	}
}
func (ctl *Controller) UpdateMenu(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[MenuInput](c)
	if ok {
		ctl.audited(c, ctl.service.UpdateMenu(c.Request.Context(), id, in), "UPDATE", "MENU", id)
	}
}
func (ctl *Controller) DeleteMenu(c *gin.Context) {
	id, ok := bindID(c)
	if ok {
		ctl.audited(c, ctl.service.DeleteMenu(c.Request.Context(), id), "DELETE", "MENU", id)
	}
}
func (ctl *Controller) SetGrants(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[GrantInput](c)
	if ok {
		ctl.audited(c, ctl.service.SetGrants(c.Request.Context(), id, in), "GRANT", "ROLE", id)
	}
}
func (ctl *Controller) SetPermissionResources(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[ResourceGrantInput](c)
	if ok {
		ctl.audited(c, ctl.service.SetPermissionResources(c.Request.Context(), id, in), "MAP", "PERMISSION", id)
	}
}
func (ctl *Controller) CreateGroup(c *gin.Context) {
	in, ok := bindJSON[PermissionGroupInput](c)
	if ok {
		ctl.audited(c, ctl.service.CreateGroup(c.Request.Context(), in), "CREATE", "PERMISSION_GROUP", 0)
	}
}
func (ctl *Controller) UpdateGroup(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[PermissionGroupInput](c)
	if ok {
		ctl.audited(c, ctl.service.UpdateGroup(c.Request.Context(), id, in), "UPDATE", "PERMISSION_GROUP", id)
	}
}
func (ctl *Controller) DeleteGroup(c *gin.Context) {
	id, ok := bindID(c)
	if ok {
		ctl.audited(c, ctl.service.DeleteGroup(c.Request.Context(), id), "DELETE", "PERMISSION_GROUP", id)
	}
}
func (ctl *Controller) CreateAPI(c *gin.Context) {
	in, ok := bindJSON[APIInput](c)
	if ok {
		ctl.audited(c, ctl.service.CreateAPI(c.Request.Context(), in), "CREATE", "API", 0)
	}
}
func (ctl *Controller) UpdateAPI(c *gin.Context) {
	id, ok := bindID(c)
	if !ok {
		return
	}
	in, ok := bindJSON[APIInput](c)
	if ok {
		ctl.audited(c, ctl.service.UpdateAPI(c.Request.Context(), id, in), "UPDATE", "API", id)
	}
}
func (ctl *Controller) DeleteAPI(c *gin.Context) {
	id, ok := bindID(c)
	if ok {
		ctl.audited(c, ctl.service.DeleteAPI(c.Request.Context(), id), "DELETE", "API", id)
	}
}
func (ctl *Controller) SyncAPIs(c *gin.Context) {
	for _, route := range ctl.engine.Routes() {
		if strings.HasPrefix(route.Path, "/swagger/") {
			continue
		}
		in := APIInput{Name: route.Handler, Method: route.Method, Path: route.Path, Enabled: boolPointer(true)}
		if err := ctl.service.repo.SyncAPI(c.Request.Context(), in); err != nil {
			ctl.audited(c, publicError(err), "SYNC", "API", 0)
			return
		}
	}
	ctl.audited(c, nil, "SYNC", "API", 0)
}
func boolPointer(value bool) *bool { return &value }
func (ctl *Controller) MyAccess(c *gin.Context) {
	id, err := utils.MustGetUser(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	out, err := ctl.service.Access(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.WriteSuccess(c, out)
}

func (ctl *Controller) Require(code string) gin.HandlerFunc {
	return ctl.RequireAny(code)
}

func (ctl *Controller) RequireAny(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.MustGetUser(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, appErrors.ErrSessionUnauthorized)
			return
		}
		allowed, err := ctl.service.AllowedAny(c.Request.Context(), id, codes...)
		if err == nil && !allowed {
			allowed, err = ctl.service.AllowedResource(c.Request.Context(), id, c.Request.Method, c.FullPath())
		}
		if err != nil {
			_ = c.AbortWithError(http.StatusForbidden, err)
			return
		}
		if !allowed {
			_ = c.AbortWithError(http.StatusForbidden, appErrors.ErrForbidden)
			return
		}
		c.Next()
	}
}
