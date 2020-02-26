package rhsm

import (
	"context"

	"github.com/umich-vci/gorhsm"
)

// Config holds the provider configuration
type Config struct {
	RefreshToken string
}

// Client returns a new client for accessing Red Hat Satellite
func (c *Config) Client() (*gorhsm.APIClient, context.Context, error) {
	config := gorhsm.NewConfiguration()

	token, err := gorhsm.GenerateAccessToken(c.RefreshToken)
	if err != nil {
		return nil, nil, err
	}

	auth := context.WithValue(context.Background(), gorhsm.ContextAPIKey, gorhsm.APIKey{
		Key:    token.AccessToken,
		Prefix: "Bearer", // Omit if not necessary.
	})

	return gorhsm.NewAPIClient(config), auth, nil
}
