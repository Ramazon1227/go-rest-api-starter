package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Ramazon1227/go-rest-api-starter/api/middleware"
	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	jwt.SigningKey = []byte("middleware-test-secret")
	m.Run()
}

func newMiddlewareRouter(mw ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(mw...)
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

// tokenFor generates a JWT for the given user. Each test uses a unique ID to
// prevent collisions with the in-memory blacklist across tests that run within
// the same second (identical claims → identical HMAC token).
func tokenFor(t *testing.T, id, role string) string {
	t.Helper()
	user := &models.User{Id: id, Email: id + "@example.com", Role: role}
	tok, err := jwt.GenerateToken(user)
	if err != nil {
		t.Fatalf("tokenFor: %v", err)
	}
	return tok
}

// ── AuthMiddleware ────────────────────────────────────────────────────────────

func TestAuthMiddleware_NoHeader(t *testing.T) {
	r := newMiddlewareRouter(middleware.AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MissingBearerPrefix(t *testing.T) {
	r := newMiddlewareRouter(middleware.AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token abc123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := newMiddlewareRouter(middleware.AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_BlacklistedToken(t *testing.T) {
	tok := tokenFor(t, "blacklist-user", "STUDENT")
	jwt.InvalidateToken(tok) //nolint:errcheck

	r := newMiddlewareRouter(middleware.AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	tok := tokenFor(t, "valid-user", "STUDENT")

	r := newMiddlewareRouter(middleware.AuthMiddleware())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ── RequireRole ───────────────────────────────────────────────────────────────

func TestRequireRole_CorrectRole(t *testing.T) {
	tok := tokenFor(t, "admin-role-user", "SYSTEM_ADMIN")
	r := newMiddlewareRouter(middleware.AuthMiddleware(), middleware.RequireRole("SYSTEM_ADMIN"))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_WrongRole(t *testing.T) {
	tok := tokenFor(t, "wrong-role-user", "STUDENT")
	r := newMiddlewareRouter(middleware.AuthMiddleware(), middleware.RequireRole("SYSTEM_ADMIN"))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_MultipleAllowedRoles(t *testing.T) {
	for _, tc := range []struct{ id, role string }{
		{"multi-admin", "SYSTEM_ADMIN"},
		{"multi-instructor", "INSTRUCTOR"},
	} {
		t.Run(tc.role, func(t *testing.T) {
			tok := tokenFor(t, tc.id, tc.role)
			r := newMiddlewareRouter(middleware.AuthMiddleware(), middleware.RequireRole("SYSTEM_ADMIN", "INSTRUCTOR"))
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tok)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
