package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/config"
	appErrors "github.com/keep/sunny/pkg/errors"
)

const (
	googleStateCookie = "google_oauth_state"
	googleNextCookie  = "google_oauth_next"
	googleStateTTL    = 10 * time.Minute
)

type Controller struct {
	service           Service
	cookieSecure      bool
	cookieName        string
	refreshCookieName string
	allowedOrigins    map[string]struct{}
}

func NewController(cfg config.Config, service Service) *Controller {
	origins := make(map[string]struct{}, len(cfg.HTTP().AllowedOrigins))
	for _, origin := range cfg.HTTP().AllowedOrigins {
		origins[origin] = struct{}{}
	}
	return &Controller{service: service, cookieSecure: cfg.HTTP().CookieSecure, cookieName: cfg.HTTP().CookieName, refreshCookieName: cfg.HTTP().RefreshCookieName, allowedOrigins: origins}
}

func newOAuthState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func safeReturnPath(value string) string {
	if value == "" || value[0] != '/' || len(value) > 1 && value[1] == '/' {
		return "/"
	}
	return value
}

func (controller *Controller) returnURL(origin, path string) string {
	path = safeReturnPath(path)
	if _, ok := controller.allowedOrigins[origin]; ok {
		return origin + path
	}
	return path
}

func (controller *Controller) safeReturnURL(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || !parsed.IsAbs() {
		return safeReturnPath(value)
	}
	origin := parsed.Scheme + "://" + parsed.Host
	if parsed.User != nil || parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "/"
	}
	if _, ok := controller.allowedOrigins[origin]; !ok {
		return "/"
	}
	return origin + safeReturnPath(parsed.RequestURI())
}

// GoogleLogin godoc
// @Summary Start Google OAuth login
// @Tags Authentication
// @Success 307
// @Failure 503 {object} errors.APIError
// @Router /api/v1/auth/google/login [get]
func (controller *Controller) GoogleLogin(ctx *gin.Context) {
	state, err := newOAuthState()
	if err != nil {
		_ = ctx.Error(appErrors.WrapAPIError(appErrors.ErrInternalServerError, err))
		return
	}
	url, err := controller.service.AuthorizationURL(state)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(googleStateCookie, state, int(googleStateTTL.Seconds()), "/api/v1/auth/google", "", controller.cookieSecure, true)
	ctx.SetCookie(googleNextCookie, controller.returnURL(ctx.Query("origin"), ctx.Query("next")), int(googleStateTTL.Seconds()), "/api/v1/auth/google", "", controller.cookieSecure, true)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary Complete Google OAuth login
// @Tags Authentication
// @Produce json
// @Param code query string true "Google authorization code"
// @Param state query string true "OAuth state"
// @Success 200 {object} LoginAPIResponse
// @Failure 400 {object} errors.APIError
// @Failure 502 {object} errors.APIError
// @Router /api/v1/auth/google/callback [get]
func (controller *Controller) GoogleCallback(ctx *gin.Context) {
	code, state := ctx.Query("code"), ctx.Query("state")
	expectedState, cookieErr := ctx.Cookie(googleStateCookie)
	next, nextErr := ctx.Cookie(googleNextCookie)
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(googleStateCookie, "", -1, "/api/v1/auth/google", "", controller.cookieSecure, true)
	ctx.SetCookie(googleNextCookie, "", -1, "/api/v1/auth/google", "", controller.cookieSecure, true)
	if code == "" || state == "" || cookieErr != nil || state != expectedState {
		_ = ctx.Error(appErrors.NewAPIError(http.StatusBadRequest, "Google 登录状态无效或已过期"))
		return
	}
	result, err := controller.service.Login(ctx.Request.Context(), code)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.SetCookie(controller.cookieName, result.AccessToken, int(time.Until(result.ExpiresAt).Seconds()), "/", "", controller.cookieSecure, true)
	ctx.SetCookie(controller.refreshCookieName, result.session.RefreshToken, int(time.Until(result.session.RefreshExpiresAt).Seconds()), "/", "", controller.cookieSecure, true)
	if nextErr != nil {
		next = "/"
	}
	ctx.Redirect(http.StatusSeeOther, controller.safeReturnURL(next))
}
