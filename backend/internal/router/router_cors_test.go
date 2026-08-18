package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/keep/sunny/config"
	"github.com/stretchr/testify/require"
)

func TestCORSPreflightAllowsConfiguredFrontendsAndKeepHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	allowedOrigins := []string{
		"https://keep.vivl.cc",
		"http://localhost:5173",
		"http://localhost:5174",
		"https://m.gotokeep.com",
		"https://m.pre.gotokeep.com",
	}
	for _, origin := range allowedOrigins {
		t.Run(origin, func(t *testing.T) {
			engine := gin.New()
			engine.Use(cors.New(newCORSConfig(config.HTTP{AllowedOrigins: allowedOrigins})))
			engine.POST("/api/v1/users", func(c *gin.Context) { c.Status(http.StatusOK) })

			request := httptest.NewRequest(http.MethodOptions, "/api/v1/users", nil)
			request.Header.Set("Origin", origin)
			request.Header.Set("Access-Control-Request-Method", http.MethodPost)
			request.Header.Set("Access-Control-Request-Headers", "content-type,authorization,x-user-id,x-version-name")
			recorder := httptest.NewRecorder()

			engine.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusNoContent, recorder.Code)
			require.Equal(t, origin, recorder.Header().Get("Access-Control-Allow-Origin"))
			require.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
			require.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "Authorization")
			require.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "X-User-Id")
			require.Contains(t, recorder.Header().Get("Access-Control-Allow-Headers"), "X-Version-Name")
		})
	}
}
