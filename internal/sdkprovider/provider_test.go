package sdkprovider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"rhsm": func() (*schema.Provider, error) {
		return New("test")(), nil
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
