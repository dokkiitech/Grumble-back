package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	firebaseauth "firebase.google.com/go/v4/auth"

	"github.com/dokkiitech/grumble-back/internal/domain/shared"
	"github.com/dokkiitech/grumble-back/internal/logging"
	"github.com/dokkiitech/grumble-back/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FirebaseTokenVerifier defines the subset of the Firebase auth client we rely on.
type FirebaseTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*firebaseauth.Token, error)
}

// AuthMiddleware validates Firebase ID tokens from the Authorization header.
type AuthMiddleware struct {
	authClient      FirebaseTokenVerifier
	authAnonymousUC *usecase.AuthAnonymousUseCase
	logger          logging.Logger
}

// NewAuthMiddleware creates a Firebase-backed authentication middleware.
func NewAuthMiddleware(
	authClient FirebaseTokenVerifier,
	authAnonymousUC *usecase.AuthAnonymousUseCase,
	logger logging.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		authClient:      authClient,
		authAnonymousUC: authAnonymousUC,
		logger:          logger,
	}
}

// Authenticate enforces authentication using the Authorization Bearer token.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawHeader := c.GetHeader("Authorization")
		tokenString, ok := extractBearerToken(rawHeader)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": "Missing or invalid bearer token"})
			c.Abort()
			return
		}

		decodedToken, err := m.authClient.VerifyIDToken(c.Request.Context(), tokenString)
		if err != nil {
			m.logger.WarnContext(c.Request.Context(), "Failed to verify Firebase token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		userUUID, err := deriveUUIDFromFirebaseUID(decodedToken.UID)
		if err != nil {
			m.logger.ErrorContext(c.Request.Context(), "Failed to derive user UUID from Firebase UID", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED", "message": "Token is missing a valid user identifier"})
			c.Abort()
			return
		}

		userID := shared.UserID(userUUID.String())

		resp, err := m.authAnonymousUC.Authenticate(c.Request.Context(), usecase.AuthAnonymousRequest{UserID: userID})
		if err != nil {
			m.logger.ErrorContext(c.Request.Context(), "Authentication failed", "user_id", userID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AUTH_FAILED", "message": "Authentication failed"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("firebase_uid", decodedToken.UID)
		c.Set("user", resp.User)
		c.Set("is_new_user", resp.IsNewUser)

		if resp.IsNewUser {
			m.logger.InfoContext(c.Request.Context(), "New user provisioned from Firebase auth", "user_id", userID, "firebase_uid", decodedToken.UID)
		}

		c.Next()
	}
}

// OptionalAuth tries to authenticate but lets the request proceed when missing or invalid.
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawHeader := c.GetHeader("Authorization")
		tokenString, ok := extractBearerToken(rawHeader)
		if !ok {
			c.Next()
			return
		}

		decodedToken, err := m.authClient.VerifyIDToken(c.Request.Context(), tokenString)
		if err != nil {
			m.logger.WarnContext(c.Request.Context(), "Optional auth token verification failed", "error", err)
			c.Next()
			return
		}

		userUUID, err := deriveUUIDFromFirebaseUID(decodedToken.UID)
		if err != nil {
			m.logger.ErrorContext(c.Request.Context(), "Optional auth failed to derive UUID", "error", err)
			c.Next()
			return
		}

		userID := shared.UserID(userUUID.String())
		resp, err := m.authAnonymousUC.Authenticate(c.Request.Context(), usecase.AuthAnonymousRequest{UserID: userID})
		if err != nil {
			m.logger.WarnContext(c.Request.Context(), "Optional auth failed", "user_id", userID, "error", err)
			c.Next()
			return
		}

		c.Set("user_id", userID)
		c.Set("firebase_uid", decodedToken.UID)
		c.Set("user", resp.User)
		c.Set("is_new_user", resp.IsNewUser)

		c.Next()
	}
}

func extractBearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}

func deriveUUIDFromFirebaseUID(uid string) (uuid.UUID, error) {
	if strings.TrimSpace(uid) == "" {
		return uuid.UUID{}, errors.New("firebase UID is empty")
	}

	// Derive a deterministic UUID so we can reuse the existing UUID-based schema.
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(uid)), nil
}
