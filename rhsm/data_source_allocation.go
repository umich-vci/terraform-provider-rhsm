package rhsm

import (
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/umich-vci/gorhsm"
)

func dataSourceAllocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAllocationRead,
		Schema: map[string]*schema.Schema{
			"uuid": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"entitlements_attached_quantity": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"entitlement_reason": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"entitlement_valid": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"entitlements": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entitlement_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sku": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAllocationRead(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	uuid := d.Get("uuid").(string)

	optional.NewString("entitlements")

	opts := &gorhsm.ShowAllocationOpts{
		Include: optional.NewString("entitlements"),
	}

	alloc, _, err := client.AllocationApi.ShowAllocation(auth, uuid, opts)
	if err != nil {
		return err
	}

	d.SetId(alloc.Body.Uuid)
	d.Set("name", alloc.Body.Name)
	d.Set("type", alloc.Body.Type)
	d.Set("version", alloc.Body.Version)
	d.Set("created_date", alloc.Body.CreatedDate.Format("2006-01-02T15:04:05.000Z"))
	d.Set("created_by", alloc.Body.CreatedBy)
	d.Set("last_modified", alloc.Body.LastModified.Format("2006-01-02T15:04:05.000Z"))
	d.Set("entitlements_attached_quantity", alloc.Body.EntitlementsAttachedQuantity)
	d.Set("entitlement_reason", alloc.Body.EntitlementsAttached.Reason)
	d.Set("entitlement_valid", alloc.Body.EntitlementsAttached.Valid)

	entitlements := []map[string]interface{}{}
	for _, x := range alloc.Body.EntitlementsAttached.Value {
		entitlement := make(map[string]interface{})
		entitlement["contract_number"] = x.ContractNumber
		entitlement["entitlement_quantity"] = x.EntitlementQuantity
		entitlement["id"] = x.Id
		entitlement["sku"] = x.Sku
		entitlements = append(entitlements, entitlement)

	}
	d.Set("entitlements", entitlements)

	return nil
}
