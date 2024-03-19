package sdkprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAllocationEntitlement() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage entitlements to a RHSM Subscription Allocation for a Red Hat Satellite server.",

		CreateContext: resourceAllocationEntitlementCreate,
		ReadContext:   resourceAllocationEntitlementRead,
		UpdateContext: resourceAllocationEntitlementUpdate,
		DeleteContext: resourceAllocationEntitlementDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"pool": {
				Description: "The ID of the pool you would like to create the entitlement from.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"quantity": {
				Description:  "The number of entitlements you would like add to the allocation/use from the pool.",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"allocation_uuid": {
				Description:  "The UUID of the subscription allocation to create the entitlement on.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"contract_number": {
				Description: "The support contract associated with the entitlement.",
				Type:        schema.TypeString,
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

func resourceAllocationEntitlementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Id()
	include := "entitlements"

	alloc, resp, err := client.AllocationAPI.ShowAllocation(auth, allocationUUID).Include(include).Execute()
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	entitlementFound := false
	for _, x := range alloc.Body.EntitlementsAttached.GetValue() {
		if x.GetId() == entitlementID {
			entitlementFound = true
			d.Set("contract_number", x.GetContractNumber())
			d.Set("quantity", x.GetEntitlementQuantity())
			d.Set("sku", x.GetSku())

		}
	}

	if !entitlementFound {
		d.SetId("")
	}

	return nil
}

func resourceAllocationEntitlementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	pool := d.Get("pool").(string)
	allocationUUID := d.Get("allocation_uuid").(string)

	pools, _, err := client.AllocationAPI.ListAllocationPools(auth, allocationUUID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	poolFound := false
	var contractNumber string
	var sku string
	for _, x := range pools.GetBody() {
		if x.GetId() == pool {
			poolFound = true
			contractNumber = x.GetContractNumber()
			sku = x.GetSku()
		}
	}
	if !poolFound {
		return diag.Errorf("Allocation %s does not have pool with id %s", allocationUUID, pool)
	}

	quantity := int32(d.Get("quantity").(int))

	alloc, _, err := client.AllocationAPI.AttachEntitlementAllocation(auth, allocationUUID).Quantity(quantity).Pool(pool).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	entitlementFound := false
	for _, x := range alloc.Body.EntitlementsAttached.GetValue() {
		if x.GetContractNumber() == contractNumber && x.GetSku() == sku {
			entitlementFound = true
			d.SetId(x.GetId())
		}
	}
	if !entitlementFound {
		return diag.Errorf("Unable to find entitlement that was created")
	}

	return resourceAllocationEntitlementRead(ctx, d, meta)
}

func resourceAllocationEntitlementUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Id()
	quantity := int32(d.Get("quantity").(int))

	_, _, err := client.AllocationAPI.UpdateEntitlementAllocation(auth, allocationUUID, entitlementID).Quantity(quantity).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAllocationEntitlementRead(ctx, d, meta)
}

func resourceAllocationEntitlementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)
	entitlementID := d.Id()

	_, err := client.AllocationAPI.RemoveAllocationEntitlement(auth, allocationUUID, entitlementID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
