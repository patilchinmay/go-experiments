package auth

import (
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type Scope string
type Role string
type Group string

type ValidatorIntf interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
	ExtractUserInfo(token *jwt.Token) (UserInfo, error)
	ExtractScopes(token *jwt.Token) ([]Scope, []Role, []Group, error)
}

var _ ValidatorIntf = &externalJWTValidator{}

// KeycloakConfig holds the configuration required to connect to a Keycloak server,
// including the base URL of the Keycloak instance, the realm to use, and the client ID
// for the application.
// nolint: lll
type KeycloakConfig struct {
	BaseURL  string `default:"http://localhost:8080" envconfig:"KEYCLOAK_BASE_URL"`
	Realm    string `default:"myrealm" envconfig:"KEYCLOAK_REALM"`
	ClientID string `default:"my-golang-app" envconfig:"KEYCLOAK_CLIENT_ID"`
}

var kcv *externalJWTValidator

// NewKeycloakValidator creates and returns a new instance of externalJWTValidator configured for Keycloak.
// The validator can be used in the rest of the application to validate JWTs issued by the specified Keycloak realm.
func NewKeycloakValidator(logger *slog.Logger, kc KeycloakConfig) *externalJWTValidator {
	issuer := fmt.Sprintf("%s/realms/%s", kc.BaseURL, kc.Realm)
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", kc.BaseURL, kc.Realm)

	if kcv == nil {
		kcv = &externalJWTValidator{
			jwksClient:              NewJWKSClient(logger, jwksURL),
			expectedIssuer:          issuer,
			expectedAuthorizedParty: kc.ClientID,
			logger:                  logger,
		}
	}

	logger.Info("externalJWTValidator for keycloak initialized")

	return kcv
}
