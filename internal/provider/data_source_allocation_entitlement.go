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
				Description:  "The UUID of the subscription allocation to create the entitlement on.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"contract_number": {
				Description: "The support contract associated with the entitlement.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"quantity": {
				Description: "The number of entitlements available in the pool.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"sku": {
				Description: "The SKU of the entitlement.",
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
	for _, x := range *alloc.Body.EntitlementsAttached.Value {
		if *x.Id == entitlementID {
			entitlementFound = true
			d.Set("contract_number", *x.ContractNumber)
			d.Set("quantity", *x.EntitlementQuantity)
			d.Set("sku", *x.Sku)

		}
	}

	if !entitlementFound {
		return diag.Errorf("Allocation %s does not have an entitlement with id %s", allocationUUID, entitlementID)
	}

	return nil
}
