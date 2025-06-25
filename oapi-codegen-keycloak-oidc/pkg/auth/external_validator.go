package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// externalJWTValidator is responsible for validating JWT tokens using an external JWKS (JSON Web Key Set) provider.
// It verifies the token's signature and ensures that the issuer and audience match the expected values.
// Fields:
//   - jwksClient:        A client used to fetch and cache JWKS keys for signature verification.
//   - expectedIssuer:    The issuer string that the JWT token must contain.
//   - expectedAuthorizedParty:  Authorized party (the party to which this token was issued).
//
// Even though this validator is designed for Keycloak,
// it can be used with any OpenID Connect compliant identity provider such as Auth0, Okta, etc.
type externalJWTValidator struct {
	jwksClient              *JWKSClient
	expectedIssuer          string
	expectedAuthorizedParty string
	logger                  *slog.Logger
}

// ValidateToken parses and validates a JWT token string using the configured JWKS client.
// It verifies the token's signing method, retrieves the appropriate public key based on the key ID (kid)
// from the token header, and validates the token's claims. Returns the parsed token if valid, or an error otherwise.
func (v *externalJWTValidator) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get key ID from token header
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("token header missing key ID")
		}

		// Fetch public key from JWKS
		return v.jwksClient.GetPublicKey(keyID)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate token claims
	if err := v.validateClaims(token); err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	return token, nil
}

// UserInfo represents extracted user information from a JWT token, including subject, email, name,
// username, and roles.
type UserInfo struct {
	Subject  string   `json:"sub"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// ExtractUserInfo extracts user information from a given JWT token's claims.
// It populates a UserInfo struct with standard fields such as Subject, Email, Name,
// and Username (from "preferred_username"). Additionally, it extracts user roles
// from the "resource_access.my-golang-app.roles" claim, which contains the application-specific
// roles for the Keycloak client.
// Returns the populated UserInfo struct or an error if the token claims are invalid.
func (v *externalJWTValidator) ExtractUserInfo(token *jwt.Token) (UserInfo, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserInfo{}, errors.New("invalid token claims")
	}

	userInfo := UserInfo{}

	if sub, ok := claims["sub"].(string); ok {
		userInfo.Subject = sub
	}

	if email, ok := claims["email"].(string); ok {
		userInfo.Email = email
	}

	if name, ok := claims["name"].(string); ok {
		userInfo.Name = name
	}

	if preferredUsername, ok := claims["preferred_username"].(string); ok {
		userInfo.Username = preferredUsername
	}

	// Extract roles from resource_access for the specific client (my-golang-app)
	if resourceAccess, ok := claims["resource_access"].(map[string]any); ok {
		if clientAccess, ok := resourceAccess[v.expectedAuthorizedParty].(map[string]any); ok {
			if roles, ok := clientAccess["roles"].([]any); ok {
				for _, role := range roles {
					if roleStr, ok := role.(string); ok {
						userInfo.Roles = append(userInfo.Roles, roleStr)
					}
				}
			}
		}
	}

	return userInfo, nil
}

// ExtractScopes extracts scopes from a JWT token for authorization purposes.
// It looks for scopes in multiple places:
// 1. "scope" claim (space-separated string)
// 2. "resource_access.my-golang-app.roles" (application-specific roles that can act as scopes)
// 3. "groups" claim (additional group-based permissions)
// Returns a combined list of all available scopes/permissions.
func (v *externalJWTValidator) ExtractScopes(token *jwt.Token) ([]Scope, []Role, []Group, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, nil, errors.New("invalid token claims")
	}

	var (
		scopes []Scope
		roles  []Role
		groups []Group
	)

	// Extract from "scope" claim (space-separated string)
	if scopeStr, ok := claims["scope"].(string); ok {
		scopeParts := strings.Split(scopeStr, " ")
		for _, scope := range scopeParts {
			if scope != "" {
				scopes = append(scopes, Scope(scope))
			}
		}
	}

	// Extract from resource_access for the specific client (my-golang-app)
	if resourceAccess, ok := claims["resource_access"].(map[string]any); ok {
		if clientAccess, ok := resourceAccess[v.expectedAuthorizedParty].(map[string]any); ok {
			if clientRoles, ok := clientAccess["roles"].([]any); ok {
				for _, role := range clientRoles {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, Role(roleStr))
					}
				}
			}
		}
	}

	// Extract from groups claim (if present)
	if groupsClaim, ok := claims["groups"].([]any); ok {
		for _, group := range groupsClaim {
			if groupStr, ok := group.(string); ok {
				groups = append(groups, Group(groupStr))
			}
		}
	}

	return scopes, roles, groups, nil
}

// validateClaims validates the standard JWT claims in the provided token.
// It checks for the presence and correctness of the following claims:
//   - "exp" (expiration time): Ensures the token has not expired.
//   - "nbf" (not before, optional): Ensures the token is not used before its valid time.
//   - "iat" (issued at, optional): Ensures the token is not issued in the future.
//   - "iss" (issuer): Ensures the token was issued by the expected issuer.
//   - "aud" (audience): Ensures the token is intended for the expected audience.
//
// Returns an error if any of the validations fail, or nil if all claims are valid.
// nolint: gocognit,funlen
func (v *externalJWTValidator) validateClaims(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims format")
	}

	// Validate expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return errors.New("token has expired")
		}
	} else {
		return errors.New("token missing expiration claim")
	}

	// Validate not before
	if nbf, ok := claims["nbf"].(float64); ok {
		if time.Now().Unix() < int64(nbf) {
			return errors.New("token not yet valid")
		}
	}

	// Validate issued at
	if iat, ok := claims["iat"].(float64); ok {
		if time.Now().Unix() < int64(iat) {
			return errors.New("token issued in the future")
		}
	}

	// Validate issuer
	if iss, ok := claims["iss"].(string); ok {
		if iss != v.expectedIssuer {
			return fmt.Errorf("invalid issuer: expected %s, got %s", v.expectedIssuer, iss)
		}
	} else {
		return errors.New("token missing issuer claim")
	}

	// Validate authorized party (azp)
	// This is the client ID for which the token was issued.
	// It is used to ensure that the token is intended for the correct application.
	if azp, ok := claims["azp"]; ok {
		switch azpValue := azp.(type) {
		case string:
			if azpValue != v.expectedAuthorizedParty {
				return fmt.Errorf("invalid authorized party: expected %s, got %s", v.expectedAuthorizedParty, azpValue)
			}
		case []any:
			found := false

			for _, a := range azpValue {
				if azpStr, ok := a.(string); ok && azpStr == v.expectedAuthorizedParty {
					found = true

					break
				}
			}

			if !found {
				return fmt.Errorf("authorized party %s not found in token", v.expectedAuthorizedParty)
			}
		default:
			return errors.New("invalid authorized party claim format")
		}
	} else {
		return errors.New("token missing authorized party claim")
	}

	return nil
}
