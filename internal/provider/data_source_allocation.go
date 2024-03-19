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
			"created_by": {
				Description: "The user account used to create the subscription allocation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_date": {
				Description: "The date and time the subscription allocation was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"entitlements_attached": {
				Description: "A list of entitlements attached to the subscription allocation.",
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
						"values": {
							Description: "A list of entitlements attached to the subscription allocation.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
									"id": {
										Description: "The ID of the entitlement.",
										Type:        schema.TypeString,
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
							},
						},
					},
				},
			},
			"entitlements_attached_quantity": {
				Description: "The number of entitlements associated with the subscription",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"last_modified": {
				Description: "The date and time the subscription allocation was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the subscription allocation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"simple_content_access": {
				Description: "Simple content access status.",
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
		},
	}
}

func dataSourceAllocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	uuid := d.Get("uuid").(string)
	include := "entitlements"

	alloc, _, err := client.AllocationAPI.ShowAllocation(auth, uuid).Include(include).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(alloc.Body.GetUuid())
	d.Set("created_by", alloc.Body.GetCreatedBy())
	d.Set("created_date", alloc.Body.GetCreatedDate())
	d.Set("entitlements_attached_quantity", alloc.Body.GetEntitlementsAttachedQuantity())
	d.Set("last_modified", alloc.Body.GetLastModified())
	d.Set("name", alloc.Body.GetName())
	d.Set("simple_content_access", alloc.Body.GetSimpleContentAccess())
	d.Set("type", alloc.Body.GetType())
	d.Set("uuid", alloc.Body.GetUuid())
	d.Set("version", alloc.Body.GetVersion())

	valueList := make([]map[string]interface{}, 0)
	for _, x := range alloc.Body.EntitlementsAttached.GetValue() {
		value := make(map[string]interface{})
		value["contract_number"] = x.GetContractNumber()
		value["end_date"] = x.GetEndDate()
		value["entitlement_quantity"] = x.GetEntitlementQuantity()
		value["id"] = x.GetId()
		value["sku"] = x.GetSku()
		value["start_date"] = x.GetStartDate()
		value["subscription_name"] = x.GetSubscriptionName()
		valueList = append(valueList, value)
	}

	entitlementsAttached := make(map[string]interface{})
	entitlementsAttached["reason"] = alloc.Body.EntitlementsAttached.GetReason()
	entitlementsAttached["valid"] = alloc.Body.EntitlementsAttached.GetValid()
	entitlementsAttached["value"] = valueList
	entitlementsAttachedList := []map[string]interface{}{entitlementsAttached}
	d.Set("entitlements_attached", entitlementsAttachedList)

	return nil
}
