package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudAccess() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to look up information about cloud providers entitled to Red Hat Cloud Access.",

		ReadContext: dataSourceCloudAccessRead,

		Schema: map[string]*schema.Schema{
			"enabled_accounts": {
				Description: "A list where each entry is a single cloud provider",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accounts": {
							Description: "A list of cloud accounts that are enabled for cloud access in the cloud provider.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The id of the cloud account.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"date_added": {
										Description: "The date the account was added to cloud access.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"gold_image_status": {
										Description: "The status of any requests for gold image access for a cloud account.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"nickname": {
										Description: "A nickname associated with the cloud account.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"name": {
							Description: "The name of the cloud provider.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"products": {
							Description: "A list of products that are entitled to the cloud provider.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled_quantity": {
										Description: "The quantity of subscriptions allowed to be consumed by the cloud provider.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"image_groups": {
										Description: "A list of images associated with the cloud provider. These are used when requesting access to gold images for a cloud account.",
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        schema.TypeString,
									},
									"name": {
										Description: "The name of the product.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"next_renewal": {
										Description: "The renewal date of the subscription.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"sku": {
										Description: "The SKU of the product.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"total_quantity": {
										Description: "The total number of subscriptions of the product available.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
								},
							},
						},
						"short_name": {
							Description: "An abreviation of the cloud provider name. Used when adding or removing accounts.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	cap, _, err := client.CloudaccessApi.ListEnabledCloudAccessProviders(auth).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudProviders := make([]map[string]interface{}, 0)

	for _, x := range *cap.Body {
		cloudProvider := make(map[string]interface{})
		cloudProvider["name"] = *x.Name
		cloudProvider["short_name"] = *x.ShortName

		accounts := make([]map[string]interface{}, 0)
		for _, y := range *x.Accounts {
			account := make(map[string]interface{})
			account["id"] = y.Id
			account["nickname"] = y.Nickname
			account["date_added"] = y.DateAdded
			account["gold_image_status"] = y.GoldImageStatus
			accounts = append(accounts, account)
		}
		cloudProvider["accounts"] = accounts

		products := make([]map[string]interface{}, 0)
		for _, y := range *x.Products {
			product := make(map[string]interface{})
			product["name"] = y.Name
			product["sku"] = y.Sku
			product["enabled_quantity"] = y.EnabledQuantity
			product["total_quantity"] = y.TotalQuantity
			product["image_groups"] = y.ImageGroups
			product["next_renewal"] = y.NextRenewal
			products = append(products, product)
		}
		cloudProvider["products"] = products

		cloudProviders = append(cloudProviders, cloudProvider)
	}

	d.SetId(time.Now().UTC().String())

	d.Set("enabled_accounts", cloudProviders)

	return nil
}
