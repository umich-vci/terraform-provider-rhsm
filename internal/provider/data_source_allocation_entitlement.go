package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAllocationEntitlement() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to get information about an entitlement associated with a Red Hat Subscription Manager allocation.",

		ReadContext: dataSourceAllocationEntitlementRead,

		Schema: map[string]*schema.Schema{
			"entitlement_id": {
				Description: "The ID of the entitlement to look up in the specified allocation.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"allocation_uuid": {
				Description:  "The UUID of the subscription allocation containing the entitlement.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"contract_number": {
				Description: "The subscription contract number.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"end_date": {
				Description: "The date the subscription ends.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"entitlement_quantity": {
				Description: "The quantity of the subscription available in the entitlement.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"sku": {
				Description: "The SKU of the subscription.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"start_date": {
				Description: "The date the subscription starts.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subscription_name": {
				Description: "The name of the subscription.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAllocationEntitlementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Get("entitlement_id").(string)
	include := "entitlements"

	alloc, _, err := client.AllocationApi.ShowAllocation(auth, allocationUUID).Include(include).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(entitlementID)

	entitlementFound := false
	for _, x := range alloc.Body.EntitlementsAttached.GetValue() {
		if x.GetId() == entitlementID {
			entitlementFound = true
			d.Set("contract_number", x.GetContractNumber())
			d.Set("end_date", x.GetEndDate())
			d.Set("entitlement_quantity", x.GetEntitlementQuantity())
			d.Set("sku", x.GetSku())
			d.Set("start_date", x.GetStartDate())
			d.Set("subscription_name", x.GetSubscriptionName())
		}
	}

	if !entitlementFound {
		return diag.Errorf("Allocation %s does not have an entitlement with id %s", allocationUUID, entitlementID)
	}

	return nil
}
