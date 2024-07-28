package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	Provider     string `envconfig:"PROVIDER" required:"true"`
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
}

func main() {
	// Load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: ", err.Error())
	}

	// Init config
	cfg := &Config{}

	// Use env vars to populate config
	err = envconfig.Process("", cfg)
	if err != nil {
		log.Fatal("error processing config with env vars: ", err.Error())
	}

	// Set up test server
	// It checks the auth header
	// Parses the token with validation
	// And returns a success response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			// parse token
			// Skip validation, since this is just for demo purpose. Do NOT use this in production.
			token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				log.Fatal("failed to parse token: ", err.Error())
			}

			// Get claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("invalid token")
			}

			// Extract the token ID
			// Since the token is issued by keycloak, the token id field is `jti``
			tokenId, ok := claims["jti"].(string)
			if !ok {
				log.Println("token ID not found")
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenId))
		} else {
			log.Println("no token in authorization header")
			w.WriteHeader(http.StatusBadGateway)
		}
	}))
	defer ts.Close()

	// Prepare OIDC provider
	ctx := context.TODO()
	provider, err := oidc.NewProvider(ctx, cfg.Provider)
	if err != nil {
		log.Fatal("failed to create OIDC provider: ", err.Error())
	}

	// Get auth and token endpoints
	endpoint := provider.Endpoint()

	// Prepare for client credential flow
	cc := clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     endpoint.TokenURL,
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "email", "profile"},
	}

	// HTTP client using the provided token. The token will auto-refresh as necessary.
	client := cc.Client(ctx)

	// loop requests
	// a better way to do this would be via for-select instead of using time.Sleep
	for i := 0; i < 10; i++ {
		// Make the request
		resp, err := client.Get(ts.URL)
		if err != nil {
			log.Fatal("error in api response: ", err.Error())
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response:", err)
		}

		log.Printf("req no: %d response: %d tokenID: %s\n", i, resp.StatusCode, body)

		time.Sleep(30 * time.Second)
	}
}
