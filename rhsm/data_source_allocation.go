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
			"entitlements_attached": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid": {
							Type:     schema.TypeBool,
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

	entitlementsAttached := make(map[string]interface{})
	entitlementsAttached["reason"] = alloc.Body.EntitlementsAttached.Reason
	entitlementsAttached["valid"] = alloc.Body.EntitlementsAttached.Valid
	entitlementsAttachedList := []map[string]interface{}{entitlementsAttached}
	d.Set("entitlements_attached", entitlementsAttachedList)

	return nil
}
