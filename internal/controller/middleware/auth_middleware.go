package middleware

import (
	"net/http"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware handles anonymous user authentication
type AuthMiddleware struct {
	authAnonymousUC *usecase.AuthAnonymousUseCase
	logger          logging.Logger
}

const (
	userIDCookieName   = "grumble_user_id"
	userIDCookieMaxAge = 365 * 24 * 60 * 60 // 1 year
)

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(
	authAnonymousUC *usecase.AuthAnonymousUseCase,
	logger logging.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		authAnonymousUC: authAnonymousUC,
		logger:          logger,
	}
}

// Authenticate resolves the user from cookie and creates one if absent
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Resolve user ID from cookie, generating a fresh one if needed
		userUUID, fromCookie := m.resolveUserUUID(c)

		userID := shared.UserID(userUUID.String())

		// Authenticate or create user
		req := usecase.AuthAnonymousRequest{
			UserID: userID,
		}

		resp, err := m.authAnonymousUC.Authenticate(c.Request.Context(), req)
		if err != nil {
			m.logger.ErrorContext(c.Request.Context(), "Authentication failed", "user_id", userID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AUTH_FAILED", "message": "Authentication failed"})
			c.Abort()
			return
		}

		// Ensure client keeps user identifier
		if !fromCookie || resp.IsNewUser {
			c.SetCookie(userIDCookieName, userUUID.String(), userIDCookieMaxAge, "/", "", c.Request.TLS != nil, true)
		}

		// Store user ID and user object in context
		c.Set("user_id", userID)
		c.Set("user", resp.User)
		c.Set("is_new_user", resp.IsNewUser)

		// Log new user creation
		if resp.IsNewUser {
			m.logger.InfoContext(c.Request.Context(), "New anonymous user created", "user_id", userID)
		}

		c.Next()
	}
}

// OptionalAuth is a middleware that authenticates the user if we have a cookie
// but allows the request to continue without authentication
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, exists := m.getUserUUIDFromCookie(c)
		if !exists {
			// No user cookie, continue without user context
			c.Next()
			return
		}

		userID := shared.UserID(userUUID.String())
		req := usecase.AuthAnonymousRequest{
			UserID: userID,
		}

		resp, err := m.authAnonymousUC.Authenticate(c.Request.Context(), req)
		if err != nil {
			m.logger.WarnContext(c.Request.Context(), "Optional authentication failed", "user_id", userID, "error", err)
			c.Next()
			return
		}

		// Store in context
		c.Set("user_id", userID)
		c.Set("user", resp.User)
		c.Set("is_new_user", resp.IsNewUser)

		c.Next()
	}
}

// resolveUserUUID returns a valid UUID for the requester.
// When absent or invalid, a new one is generated.
func (m *AuthMiddleware) resolveUserUUID(c *gin.Context) (uuid.UUID, bool) {
	if cookieUUID, ok := m.getUserUUIDFromCookie(c); ok {
		return cookieUUID, true
	}

	newUUID := uuid.New()
	m.logger.InfoContext(c.Request.Context(), "Generated new anonymous user ID", "user_id", newUUID.String())
	return newUUID, false
}

func (m *AuthMiddleware) getUserUUIDFromCookie(c *gin.Context) (uuid.UUID, bool) {
	cookieVal, err := c.Cookie(userIDCookieName)
	if err != nil || cookieVal == "" {
		return uuid.UUID{}, false
	}

	userUUID, err := uuid.Parse(cookieVal)
	if err != nil {
		m.logger.WarnContext(c.Request.Context(), "Invalid user ID cookie", "value", cookieVal, "error", err)
		return uuid.UUID{}, false
	}
	return userUUID, true
}
