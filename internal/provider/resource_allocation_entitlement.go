package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAllocationEntitlement() *schema.Resource {
	return &schema.Resource{
		Create: resourceAllocationEntitlementCreate,
		Read:   resourceAllocationEntitlementRead,
		Update: resourceAllocationEntitlementUpdate,
		Delete: resourceAllocationEntitlementDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quantity": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"allocation_uuid": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"contract_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAllocationEntitlementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Id()
	include := "entitlements"

	alloc, resp, err := client.AllocationApi.ShowAllocation(auth, allocationUUID).Include(include).Execute()
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

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
		d.SetId("")
	}

	return nil
}

func resourceAllocationEntitlementCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	pool := d.Get("pool").(string)
	allocationUUID := d.Get("allocation_uuid").(string)

	pools, _, err := client.AllocationApi.ListAllocationPools(auth, allocationUUID).Execute()
	if err != nil {
		return err
	}

	poolFound := false
	var contractNumber string
	var sku string
	for _, x := range *pools.Body {
		if *x.Id == pool {
			poolFound = true
			contractNumber = *x.ContractNumber
			sku = *x.Sku
		}
	}
	if !poolFound {
		return fmt.Errorf("Allocation %s does not have pool with id %s", allocationUUID, pool)
	}

	quantity := int32(d.Get("quantity").(int))

	alloc, _, err := client.AllocationApi.AttachEntitlementAllocation(auth, allocationUUID).Quantity(quantity).Pool(pool).Execute()
	if err != nil {
		return err
	}

	entitlementFound := false
	for _, x := range *alloc.Body.EntitlementsAttached.Value {
		if *x.ContractNumber == contractNumber && *x.Sku == sku {
			entitlementFound = true
			d.SetId(*x.Id)
		}
	}
	if !entitlementFound {
		return fmt.Errorf("Unable to find entitlement that was created")
	}

	return resourceAllocationEntitlementRead(d, meta)
}

func resourceAllocationEntitlementUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Id()
	quantity := int32(d.Get("quantity").(int))

	_, _, err := client.AllocationApi.UpdateEntitlementAllocation(auth, allocationUUID, entitlementID).Quantity(quantity).Execute()
	if err != nil {
		return err
	}

	return resourceAllocationEntitlementRead(d, meta)
}

func resourceAllocationEntitlementDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Id()

	_, err := client.AllocationApi.RemoveAllocationEntitlement(auth, allocationUUID, entitlementID).Execute()
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
