package middleware

import (
	"context"
	"strings"

	"github.com/Energie-Burgenland/ausaestung-info/utils/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware validates the incoming JWT token from Microsoft Entra ID.
func JWTMiddleware(verifier *oidc.IDTokenVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the Authorization header
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid authorization header",
			})
		}

		// Extract the token string
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the token
		idToken, err := verifier.Verify(context.Background(), rawToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Parse claims if needed
		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse claims",
			})
		}

		userName, exists := claims["preferred_username"].(string)
		if !exists {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "preferred_username claim not found",
			})
		}

		roles, exists := claims["roles"].([]interface{})
		if !exists {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "roles claim not found",
			})
		}

		role := "none"
		if len(roles) > 0 {
			role = roles[0].(string)
		}

		auth.SetUserName(userName)
		auth.SetRole(role)

		return c.Next()
	}
}
