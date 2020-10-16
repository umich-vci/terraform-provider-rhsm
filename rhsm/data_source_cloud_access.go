package rhsm

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudAccessRead,
		Schema: map[string]*schema.Schema{
			"enabled_accounts": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"date_added": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"gold_image_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nickname": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"products": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled_quantity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image_groups": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     schema.TypeString,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"next_renewal": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sku": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"total_quantity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"short_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudAccessRead(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	cap, _, err := client.CloudaccessApi.ListEnabledCloudAccessProviders(auth)
	if err != nil {
		return err
	}

	cloudProviders := make([]map[string]interface{}, 0)

	for _, x := range cap.Body {
		cloudProvider := make(map[string]interface{})
		cloudProvider["name"] = x.Name
		cloudProvider["short_name"] = x.ShortName

		accounts := make([]map[string]interface{}, 0)
		for _, y := range x.Accounts {
			account := make(map[string]interface{})
			account["id"] = y.Id
			account["nickname"] = y.Nickname
			account["date_added"] = y.DateAdded
			account["gold_image_status"] = y.GoldImageStatus
			accounts = append(accounts, account)
		}
		cloudProvider["accounts"] = accounts

		products := make([]map[string]interface{}, 0)
		for _, y := range x.Products {
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
