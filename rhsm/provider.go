package rhsm

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		RefreshToken: d.Get("refresh_token").(string),
	}

	return config, nil
}

// Provider returns a terraform resource provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RHSM_REFRESH_TOKEN", nil),
				Description: "RHSM API Refresh Token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rhsm_allocation":             resourceAllocation(),
			"rhsm_allocation_entitlement": resourceAllocationEntitlement(),
			"rhsm_allocation_manifest":    resourceAllocationManifest(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rhsm_allocation":             dataSourceAllocation(),
			"rhsm_allocation_entitlement": dataSourceAllocationEntitlement(),
			"rhsm_allocation_pools":       dataSourceAllocationPools(),
		},
		ConfigureFunc: providerConfigure,
	}
}
