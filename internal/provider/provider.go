package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/umich-vci/gorhsm"
)

// Ensure RHSMProvider satisfies various provider interfaces.
var _ provider.Provider = &RHSMProvider{}

type RHSMProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RHSMProviderModel describes the provider data model.
type RHSMProviderModel struct {
	RefreshToken types.String `tfsdk:"refresh_token"`
}

type apiClient struct {
	Auth   context.Context
	Client *gorhsm.APIClient
}

func (p *RHSMProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "rhsm"
	resp.Version = p.version
}

func (p *RHSMProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "This is the [offline token](https://access.redhat.com/articles/3626371#bgenerating-a-new-offline-tokenb-3) used to generate access tokens for Red Hat Subscription Manager. This must be provided in the config or in the environment variable `RHSM_REFRESH_TOKEN`.",
				Optional:            true,
			},
		},
	}
}

func (p *RHSMProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config RHSMProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.RefreshToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Missing refresh_token",
			"The provider cannot create the RHSM client as there is an unknown configuration value for the API refresh_token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RHSM_REFRESH_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	refreshToken := os.Getenv("RHSM_REFRESH_TOKEN")

	if !config.RefreshToken.IsNull() {
		refreshToken = config.RefreshToken.ValueString()
	}

	if refreshToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Missing refresh_token",
			"The provider cannot create the RHSM client as there is a missing or empty value for the refresh_token. "+
				"Set the value in the configuration or use the RHSM_REFRESH_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	rhsmConfig := gorhsm.NewConfiguration()
	token, err := gorhsm.GenerateAccessToken(refreshToken)
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate access token", err.Error())
		return
	}

	tokenMap := map[string]gorhsm.APIKey{"Bearer": {
		Key:    token.AccessToken,
		Prefix: token.TokenType,
	}}

	rhsmClient := &apiClient{
		Auth:   context.WithValue(context.Background(), gorhsm.ContextAPIKeys, tokenMap),
		Client: gorhsm.NewAPIClient(rhsmConfig),
	}

	// Make the BlueCat client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = rhsmClient
	resp.ResourceData = rhsmClient

}

func (p *RHSMProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCloudAccessAccountResource,
	}
}

func (p *RHSMProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) provider.Provider {
	return &RHSMProvider{
		version: version,
	}
}
