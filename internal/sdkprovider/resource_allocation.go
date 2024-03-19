package sdkprovider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var nameRegex, _ = regexp.Compile(`^[a-zA-Z0-9\_\-\.]{1,100}$`)

func resourceAllocation() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a RHSM Subscription allocation for a Red Hat Satellite server.",

		CreateContext: resourceAllocationCreate,
		ReadContext:   resourceAllocationRead,
		DeleteContext: resourceAllocationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the subscription allocation.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(nameRegex, "name must be less than 100 characters and can use only numbers, letters, underscores, hyphens, and periods"),
			},
			"uuid": {
				Description: "The UUID of the subscription allocation that was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The type of the subscription allocation.  The only one supported by this resource is `Satellite`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description: "The version of the subscription allocation type.  This defaults in the API to 6.5 and cannot be set through the API. It can be adjusted in the RHSM portal.",
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
				Description: "",
				Type:        schema.TypeList,
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

func resourceAllocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	uuid := d.Id()
	include := "entitlements"

	alloc, resp, err := client.AllocationAPI.ShowAllocation(auth, uuid).Include(include).Execute()
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.SetId(alloc.Body.GetUuid())
	d.Set("name", alloc.Body.GetName())
	d.Set("uuid", alloc.Body.GetUuid())
	d.Set("type", alloc.Body.GetType())
	d.Set("version", alloc.Body.GetVersion())
	d.Set("created_date", alloc.Body.GetCreatedDate())
	d.Set("created_by", alloc.Body.GetCreatedBy())
	d.Set("last_modified", alloc.Body.GetLastModified())
	d.Set("entitlements_attached_quantity", alloc.Body.GetEntitlementsAttachedQuantity())

	entitlementsAttached := make(map[string]interface{})
	entitlementsAttached["reason"] = alloc.Body.EntitlementsAttached.GetReason()
	entitlementsAttached["valid"] = alloc.Body.EntitlementsAttached.GetValid()
	entitlementsAttachedList := []map[string]interface{}{entitlementsAttached}
	d.Set("entitlements_attached", entitlementsAttachedList)

	return nil
}

func resourceAllocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	name := d.Get("name").(string)

	alloc, _, err := client.AllocationAPI.CreateSatellite(auth).Name(name).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alloc.Body.GetUuid())

	return resourceAllocationRead(ctx, d, meta)
}

func resourceAllocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	uuid := d.Id()

	_, err := client.AllocationAPI.RemoveAllocation(auth, uuid).Force(true).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
