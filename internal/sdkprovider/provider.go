package sdkprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gorhsm"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"refresh_token": {
					Description: "This is the [offline token](https://access.redhat.com/articles/3626371#bgenerating-a-new-offline-tokenb-3) used to generate access tokens for Red Hat Subscription Manager. This must be provided in the config or in the environment variable `RHSM_REFRESH_TOKEN`.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RHSM_REFRESH_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"rhsm_allocation":             dataSourceAllocation(),
				"rhsm_allocation_entitlement": dataSourceAllocationEntitlement(),
				"rhsm_allocation_pools":       dataSourceAllocationPools(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"rhsm_allocation":             resourceAllocation(),
				"rhsm_allocation_entitlement": resourceAllocationEntitlement(),
				"rhsm_allocation_manifest":    resourceAllocationManifest(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	Auth   context.Context
	Client *gorhsm.APIClient
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		userAgent := p.UserAgent("terraform-provider-rhsm", version)
		refreshToken := d.Get("refresh_token").(string)

		config := gorhsm.NewConfiguration()

		token, err := gorhsm.GenerateAccessToken(refreshToken)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		config.UserAgent = userAgent

		tokenMap := map[string]gorhsm.APIKey{"Bearer": {
			Key:    token.AccessToken,
			Prefix: token.TokenType,
		}}

		auth := context.WithValue(context.Background(), gorhsm.ContextAPIKeys, tokenMap)

		client := gorhsm.NewAPIClient(config)

		return &apiClient{Auth: auth, Client: client}, nil
	}
}
