package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/keep/sunny/config"
	"github.com/keep/sunny/ent"
	"github.com/keep/sunny/ent/role"
	"github.com/keep/sunny/ent/user"
	appErrors "github.com/keep/sunny/pkg/errors"
	"github.com/keep/sunny/pkg/passwordpolicy"
	"github.com/keep/sunny/pkg/postgres"
	"github.com/keep/sunny/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

const issuer = "referral"

type claims struct {
	Kind string `json:"kind"`
	jwt.RegisteredClaims
}

// Manager validates credentials and owns the authenticated HTTP session.
type Manager struct {
	client       *ent.Client
	http         config.HTTP
	jwt          config.JWT
	app          config.App
	registration config.Registration
}

// NewManager creates the authentication manager.
func NewManager(pg *postgres.Postgres, cfg config.Config) *Manager {
	return &Manager{client: pg.Client, http: cfg.HTTP(), jwt: cfg.JWT(), app: cfg.App(), registration: cfg.Registration()}
}

func parseTTL(value string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return fallback
	}
	return duration
}

func (m *Manager) sign(userID int, kind string, ttl time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{Kind: kind, RegisteredClaims: jwt.RegisteredClaims{Issuer: issuer, Subject: strconv.Itoa(userID), IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(ttl))}})
	return token.SignedString([]byte(m.jwt.Secret))
}

// Issue creates access and refresh tokens for an authenticated user.
func (m *Manager) Issue(account *ent.User) (SessionTokens, error) {
	accessTTL := parseTTL(m.jwt.AccessTTL, 15*time.Minute)
	refreshTTL := parseTTL(m.jwt.RefreshTTL, 7*24*time.Hour)
	accessToken, err := m.sign(account.ID, "access", accessTTL)
	if err != nil {
		return SessionTokens{}, fmt.Errorf("sign access token: %w", err)
	}
	refreshToken, err := m.sign(account.ID, "refresh", refreshTTL)
	if err != nil {
		return SessionTokens{}, fmt.Errorf("sign refresh token: %w", err)
	}
	now := time.Now()
	return SessionTokens{AccessToken: accessToken, AccessExpiresAt: now.Add(accessTTL), RefreshToken: refreshToken, RefreshExpiresAt: now.Add(refreshTTL)}, nil
}

// StartSession signs and stores a browser session for a user.
func (m *Manager) StartSession(c *gin.Context, userID int) error {
	accessTTL := parseTTL(m.jwt.AccessTTL, 15*time.Minute)
	refreshTTL := parseTTL(m.jwt.RefreshTTL, 7*24*time.Hour)
	access, err := m.sign(userID, "access", accessTTL)
	if err != nil {
		return fmt.Errorf("sign access token: %w", err)
	}
	refresh, err := m.sign(userID, "refresh", refreshTTL)
	if err != nil {
		return fmt.Errorf("sign refresh token: %w", err)
	}
	m.setCookie(c, m.http.CookieName, access, accessTTL)
	m.setCookie(c, m.http.RefreshCookieName, refresh, refreshTTL)
	return nil
}

func (m *Manager) verify(raw, kind string) (int, error) {
	parsed, err := jwt.ParseWithClaims(raw, &claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(m.jwt.Secret), nil
	}, jwt.WithIssuer(issuer))
	if err != nil || !parsed.Valid {
		return 0, appErrors.ErrSessionUnauthorized
	}
	c, ok := parsed.Claims.(*claims)
	if !ok || c.Kind != kind {
		return 0, appErrors.ErrSessionUnauthorized
	}
	id, err := strconv.Atoi(c.Subject)
	if err != nil || id < 1 {
		return 0, appErrors.ErrSessionUnauthorized
	}
	return id, nil
}

func (m *Manager) setCookie(c *gin.Context, name, value string, ttl time.Duration) {
	http.SetCookie(c.Writer, &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: int(ttl.Seconds()), HttpOnly: true, Secure: m.http.CookieSecure, SameSite: http.SameSiteLaxMode})
}

// Identity loads a verified access-token identity into the Gin context.
func (m *Manager) Identity() gin.HandlerFunc {
	return func(c *gin.Context) {
		if raw, err := c.Cookie(m.http.CookieName); err == nil {
			if id, verifyErr := m.verify(raw, "access"); verifyErr == nil {
				c.Set("userId", id)
			}
		}
		c.Next()
	}
}

// Require rejects requests without a verified authenticated identity.
func (m *Manager) Require() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get("userId"); !ok {
			_ = c.AbortWithError(http.StatusUnauthorized, appErrors.ErrSessionUnauthorized)
			return
		}
		c.Next()
	}
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type registrationCodeInput struct {
	Email string `json:"email" binding:"required,email,max=254"`
}

type registerInput struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	Email           string `json:"email" binding:"required,email,max=254"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=72"`
	Code            string `json:"code" binding:"required,len=6,numeric"`
}

func (m *Manager) demoRegistration() (string, bool) {
	if !m.registration.DemoEnabled && !m.app.IsDev() {
		return "", false
	}
	if m.registration.DemoCode != "" {
		return m.registration.DemoCode, true
	}
	return "123456", true
}

func (m *Manager) registrationCode(c *gin.Context) {
	var in registrationCodeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		_ = c.Error(appErrors.ErrBadRequest)
		return
	}
	code, ok := m.demoRegistration()
	if !ok {
		_ = c.Error(appErrors.NewAPIError(http.StatusServiceUnavailable, "注册验证码服务尚未配置"))
		return
	}
	response.WriteSuccess(c, gin.H{"demo": true, "code": code, "message": "演示环境不会发送邮件，请使用页面显示的验证码"})
}

// ValidateRegistrationCode checks the configured registration verification code.
func (m *Manager) ValidateRegistrationCode(value string) error {
	code, ok := m.demoRegistration()
	if !ok || value != code {
		return appErrors.NewAPIError(http.StatusBadRequest, "验证码错误或已过期")
	}
	return nil
}

func (m *Manager) register(c *gin.Context) {
	var in registerInput
	if err := c.ShouldBindJSON(&in); err != nil {
		_ = c.Error(appErrors.ErrBadRequest)
		return
	}
	if in.Password != in.ConfirmPassword {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, "两次输入的密码不一致"))
		return
	}
	if err := passwordpolicy.Validate(in.Password); err != nil {
		_ = c.Error(appErrors.NewAPIError(http.StatusBadRequest, err.Error()))
		return
	}
	if err := m.ValidateRegistrationCode(in.Code); err != nil {
		_ = c.Error(err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		_ = c.Error(err)
		return
	}
	account, err := m.client.User.Create().
		SetName(strings.TrimSpace(in.Name)).
		SetEmail(strings.ToLower(strings.TrimSpace(in.Email))).
		SetPasswordHash(string(hash)).
		Save(c.Request.Context())
	if ent.IsConstraintError(err) {
		_ = c.Error(appErrors.NewAPIError(http.StatusConflict, "该邮箱已注册"))
		return
	}
	if err != nil {
		_ = c.Error(err)
		return
	}
	accessTTL, refreshTTL := parseTTL(m.jwt.AccessTTL, 15*time.Minute), parseTTL(m.jwt.RefreshTTL, 7*24*time.Hour)
	access, err := m.sign(account.ID, "access", accessTTL)
	if err != nil {
		_ = c.Error(err)
		return
	}
	refresh, err := m.sign(account.ID, "refresh", refreshTTL)
	if err != nil {
		_ = c.Error(err)
		return
	}
	m.setCookie(c, m.http.CookieName, access, accessTTL)
	m.setCookie(c, m.http.RefreshCookieName, refresh, refreshTTL)
	response.WriteSuccess(c, gin.H{"user_id": account.ID, "name": account.Name})
}

func (m *Manager) login(c *gin.Context) {
	var in loginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		_ = c.Error(appErrors.ErrBadRequest)
		return
	}
	account, err := m.client.User.Query().Where(user.EmailEQ(strings.ToLower(strings.TrimSpace(in.Email))), user.EnabledEQ(true)).Only(c.Request.Context())
	if err != nil || bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(in.Password)) != nil {
		_ = c.Error(appErrors.ErrUnauthorized)
		return
	}
	tokens, err := m.Issue(account)
	if err != nil {
		_ = c.Error(err)
		return
	}
	m.setCookie(c, m.http.CookieName, tokens.AccessToken, time.Until(tokens.AccessExpiresAt))
	m.setCookie(c, m.http.RefreshCookieName, tokens.RefreshToken, time.Until(tokens.RefreshExpiresAt))
	response.WriteSuccess(c, gin.H{"user_id": account.ID, "name": account.Name})
}

func (m *Manager) refresh(c *gin.Context) {
	raw, err := c.Cookie(m.http.RefreshCookieName)
	if err != nil {
		_ = c.Error(appErrors.ErrSessionUnauthorized)
		return
	}
	id, err := m.verify(raw, "refresh")
	if err != nil {
		_ = c.Error(err)
		return
	}
	if exists, queryErr := m.client.User.Query().Where(user.IDEQ(id), user.EnabledEQ(true)).Exist(c.Request.Context()); queryErr != nil || !exists {
		_ = c.Error(appErrors.ErrSessionUnauthorized)
		return
	}
	ttl := parseTTL(m.jwt.AccessTTL, 15*time.Minute)
	access, err := m.sign(id, "access", ttl)
	if err != nil {
		_ = c.Error(err)
		return
	}
	m.setCookie(c, m.http.CookieName, access, ttl)
	response.WriteSuccess(c, gin.H{"refreshed": true})
}

func (m *Manager) logout(c *gin.Context) {
	m.setCookie(c, m.http.CookieName, "", -time.Hour)
	m.setCookie(c, m.http.RefreshCookieName, "", -time.Hour)
	response.WriteSuccess(c, gin.H{"logged_out": true})
}

func (m *Manager) session(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		_ = c.Error(appErrors.ErrSessionUnauthorized)
		return
	}
	account, err := m.client.User.Get(c.Request.Context(), id.(int))
	if err != nil {
		_ = c.Error(appErrors.ErrSessionUnauthorized)
		return
	}
	response.WriteSuccess(c, gin.H{"id": account.ID, "name": account.Name, "email": account.Email})
}

func (m *Manager) bootstrap(ctx context.Context, cfg config.Config) error {
	return m.bootstrapAdmin(ctx, cfg.Admin())
}

func (m *Manager) bootstrapAdmin(ctx context.Context, admin config.Admin) error {
	b := admin.Bootstrap
	if !b.Enabled && len(admin.UserIDs) == 0 && len(admin.Emails) == 0 {
		return nil
	}
	adminRole, err := m.client.Role.Query().Where(role.CodeEQ("super_admin")).Only(ctx)
	if err != nil {
		return err
	}
	if b.Enabled {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(b.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		email := strings.ToLower(strings.TrimSpace(b.Email))
		name := strings.TrimSpace(b.Name)
		account, queryErr := m.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
		if ent.IsNotFound(queryErr) {
			nameTaken, nameErr := m.client.User.Query().Where(user.NameEQ(name)).Exist(ctx)
			if nameErr != nil {
				return nameErr
			}
			if nameTaken {
				name = email
			}
			account, queryErr = m.client.User.Create().SetName(name).SetEmail(email).SetPasswordHash(string(hash)).Save(ctx)
		} else if queryErr == nil {
			nameTaken, nameErr := m.client.User.Query().Where(user.NameEQ(name), user.IDNEQ(account.ID)).Exist(ctx)
			if nameErr != nil {
				return nameErr
			}
			update := account.Update().SetPasswordHash(string(hash))
			if !nameTaken {
				update.SetName(name)
			}
			account, queryErr = update.Save(ctx)
		}
		if queryErr != nil {
			return queryErr
		}
		if err = account.Update().AddRoles(adminRole).Exec(ctx); err != nil {
			return err
		}
	}
	if len(admin.UserIDs) > 0 {
		if _, err = m.client.User.Update().Where(user.IDIn(admin.UserIDs...)).AddRoles(adminRole).Save(ctx); err != nil {
			return err
		}
	}
	emails := make([]string, 0, len(admin.Emails))
	for _, email := range admin.Emails {
		if normalized := strings.ToLower(strings.TrimSpace(email)); normalized != "" {
			emails = append(emails, normalized)
		}
	}
	if len(emails) > 0 {
		if _, err = m.client.User.Update().Where(user.EmailIn(emails...)).AddRoles(adminRole).Save(ctx); err != nil {
			return err
		}
	}
	return nil
}
