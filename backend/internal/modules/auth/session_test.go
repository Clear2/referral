package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/config"
)

func TestIdentityIgnoresSpoofedUserHeader(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	manager := &Manager{http: config.HTTP{CookieName: "session"}, jwt: config.JWT{Secret: "test-secret"}}
	router := gin.New()
	router.Use(manager.Identity())
	router.GET("/", func(c *gin.Context) {
		_, exists := c.Get("userId")
		c.JSON(http.StatusOK, gin.H{"authenticated": exists})
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-User-ID", "1")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Body.String() != "{\"authenticated\":false}" {
		t.Fatalf("spoofed identity response = %s", recorder.Body.String())
	}
}

func TestIdentityAcceptsSignedCookie(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	manager := &Manager{http: config.HTTP{CookieName: "session"}, jwt: config.JWT{Secret: "test-secret"}}
	raw, err := manager.sign(42, "access", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(manager.Identity())
	router.GET("/", func(c *gin.Context) { id, _ := c.Get("userId"); c.JSON(http.StatusOK, gin.H{"user_id": id}) })
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: "session", Value: raw})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Body.String() != "{\"user_id\":42}" {
		t.Fatalf("signed identity response = %s", recorder.Body.String())
	}
}
