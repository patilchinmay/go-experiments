package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/pkg/auth"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
	ScopesCtxKey   contextKey = "scopes"
	RolesCtxKey    contextKey = "roles"
	GroupsCtxKey   contextKey = "groups"
)

// AuthenticateAndAuthorize is a middleware function that validates the token and checks for required roles.
func (app *application) AuthenticateAndAuthorize(c echo.Context, requiredRoles []string) error {
	// Set Vary header for caching
	c.Response().Header().Add("Vary", "Authorization")

	// Extract Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
	}

	// Parse Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := parts[1]

	// Validate token
	token, err := app.validator.ValidateToken(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
	}

	// Extract user information
	userInfo, err := app.validator.ExtractUserInfo(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Failed to extract user info: %v", err))
	}

	// Extract scopes, roles and groups from token
	presentScopes, presentRoles, presentGroups, err := app.validator.ExtractScopes(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Failed to extract scopes, roles and groups from token: %v", err))
	}

	// Store user info, scopes, roles and groups in Echo context
	c.Set(string(UserContextKey), userInfo)
	c.Set(string(ScopesCtxKey), presentScopes)
	c.Set(string(RolesCtxKey), presentRoles)
	c.Set(string(GroupsCtxKey), presentGroups)

	// AuthZ Check
	app.logger.Debug("authz check",
		slog.Any("requiredRoles", requiredRoles),
		slog.Any("presentScopes", presentScopes),
		slog.Any("presentRoles", presentRoles),
		slog.Any("presentGroups", presentGroups),
	)

	// TODO: Scope check
	// TODO: Group check

	// Role check
	if !hasAllRequiredRoles(presentRoles, requiredRoles) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions, missing required roles")
	}

	return nil
}

// Helper function to check if user has ALL required roles.
func hasAllRequiredRoles(presentRoles []auth.Role, requiredRoles []string) bool {
	for _, required := range requiredRoles {
		found := false

		for _, presentRole := range presentRoles {
			if presentRole == auth.Role(required) {
				found = true

				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// // Helper function to check if user has at least one of the required scopes.
// func hasAnyScope(userScopes, requiredScopes []string) bool {
// 	for _, required := range requiredScopes {
// 		for _, userScope := range userScopes {
// 			if userScope == required {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }
