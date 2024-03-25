package sdkprovider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/umich-vci/terraform-provider-rhsm/internal/provider"
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"rhsm": func() (tfprotov6.ProviderServer, error) {
			upgradedSdkProvider, err := tf5to6server.UpgradeServer(
				context.Background(),
				New("test")().GRPCProvider,
			)

			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				func() tfprotov6.ProviderServer {
					return upgradedSdkProvider
				},

				providerserver.NewProtocol6(provider.New("test")),
			}
			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	refreshToken := os.Getenv("RHSM_REFRESH_TOKEN")

	if refreshToken == "" {
		t.Fatalf("RHSM_REFRESH_TOKEN must be set for acceptance tests to run")
	}
}
