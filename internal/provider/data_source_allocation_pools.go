package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAllocationPools() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to get information about pools available to a Red Hat Subscription Manager allocation.",

		ReadContext: dataSourceAllocationPoolsRead,

		Schema: map[string]*schema.Schema{
			"allocation_uuid": {
				Description:  "The UUID of the subscription allocation.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"future": {
				Description: "Should pools only valid in the future be listed?",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"pools": {
				Description: "A list of pools available to the subscription allocation.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_number": {
							Description: "The support contract associated with the entitlement.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"end_date": {
							Description: "The date the support contract ends.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"entitlements_available": {
							Description: "The number of entitlements available from the pool.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"id": {
							Description: "The ID of the pool.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"service_level": {
							Description: "The service level of the support contract.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"sku": {
							Description: "The SKU of the entitlement.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"start_date": {
							Description: "The date the support contract starts.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subscription_name": {
							Description: "The friendly name of the sku.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subscription_number": {
							Description: "The subscription number.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAllocationPoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	uuid := d.Get("allocation_uuid").(string)

	future := false
	if b, ok := d.GetOk("future"); ok {
		future = b.(bool)
	}

	pools, _, err := client.AllocationApi.ListAllocationPools(auth, uuid).Future(future).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uuid)

	poolsList := []map[string]interface{}{}
	for _, x := range *pools.Body {
		pool := make(map[string]interface{})
		pool["contract_number"] = *x.ContractNumber
		pool["end_date"] = *x.EndDate
		pool["entitlements_available"] = *x.EntitlementsAvailable
		pool["id"] = *x.Id
		pool["service_level"] = *x.ServiceLevel
		pool["sku"] = *x.Sku
		pool["start_date"] = *x.StartDate
		pool["subscription_name"] = *x.SubscriptionName
		pool["subscription_number"] = *x.SubscriptionNumber

		poolsList = append(poolsList, pool)

	}
	d.Set("pools", poolsList)

	return nil
}