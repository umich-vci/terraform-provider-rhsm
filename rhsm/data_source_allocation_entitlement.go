package rhsm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAllocationEntitlement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAllocationEntitlementRead,
		Schema: map[string]*schema.Schema{
			"entitlement_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"allocation_uuid": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"contract_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quantity": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAllocationEntitlementRead(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Get("entitlement_id").(string)
	include := "entitlements"

	alloc, _, err := client.AllocationApi.ShowAllocation(auth, allocationUUID).Include(include).Execute()
	if err != nil {
		return err
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
		return fmt.Errorf("Allocation %s does not have an entitlement with id %s", allocationUUID, entitlementID)
	}

	return nil
}
