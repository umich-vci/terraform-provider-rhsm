package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAllocation() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to look up a RHSM Subscription allocation.",

		ReadContext: dataSourceAllocationRead,

		Schema: map[string]*schema.Schema{
			"uuid": {
				Description:  "The UUID of the subscription allocation to look up.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description: "The name of the subscription allocation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the subscription allocation.  The only one supported by this resource is `Satellite`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description: "The version of the subscription allocation type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_date": {
				Description: "The date and time the subscription allocation was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_by": {
				Description: "The user account used to create the subscription allocation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified": {
				Description: "The date and time the subscription allocation was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"entitlements_attached_quantity": {
				Description: "The number of entitlements associated with the subscription",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"entitlements_attached": {
				Description: "A list of entitlements attached to the subscription allocation.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Description: "The reason for the value of `valid`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"valid": {
							Description: "If the entitlements associated with the subscription allocation are valid or not.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAllocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	uuid := d.Get("uuid").(string)
	include := "entitlements"

	alloc, _, err := client.AllocationApi.ShowAllocation(auth, uuid).Include(include).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alloc.Body.GetUuid())
	d.Set("name", alloc.Body.GetName())
	d.Set("type", alloc.Body.GetType())
	d.Set("version", alloc.Body.GetVersion())
	d.Set("created_date", alloc.Body.GetCreatedDate())
	d.Set("created_by", alloc.Body.GetCreatedBy())
	d.Set("last_modified", alloc.Body.GetLastModified())
	d.Set("entitlements_attached_quantity", alloc.Body.GetLastModified())

	entitlementsAttached := make(map[string]interface{})
	entitlementsAttached["reason"] = alloc.Body.EntitlementsAttached.GetReason()
	entitlementsAttached["valid"] = alloc.Body.EntitlementsAttached.GetValid()
	entitlementsAttachedList := []map[string]interface{}{entitlementsAttached}
	d.Set("entitlements_attached", entitlementsAttachedList)

	return nil
}
