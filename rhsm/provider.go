package rhsm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		RefreshToken: d.Get("refresh_token").(string),
	}

	return config, nil
}

// Provider returns a terraform resource provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RHSM_REFRESH_TOKEN", nil),
				Description: "RHSM API Refresh Token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"rhsm_allocation": dataSourceAllocation(),
		},
		ConfigureFunc: providerConfigure,
	}
}