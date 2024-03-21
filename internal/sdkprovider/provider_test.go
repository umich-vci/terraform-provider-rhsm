package sdkprovider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/umich-vci/terraform-provider-rhsm/internal/provider"
)

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"rhsm": func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()
		providers := []func() tfprotov5.ProviderServer{
			New("test")().GRPCProvider,
			providerserver.NewProtocol5(provider.New("test")),
		}

		muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
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
