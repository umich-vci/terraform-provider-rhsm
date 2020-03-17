package rhsm

import (
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/umich-vci/gorhsm"
)

func dataSourceAllocationPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAllocationPoolsRead,
		Schema: map[string]*schema.Schema{
			"allocation_uuid": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"future": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pools": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entitlements_available": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sku": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAllocationPoolsRead(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	uuid := d.Get("allocation_uuid").(string)

	opts := &gorhsm.ListAllocationPoolsOpts{}

	if b, ok := d.GetOk("future"); ok {
		opts.Future = optional.NewBool(b.(bool))
	}

	pools, _, err := client.AllocationApi.ListAllocationPools(auth, uuid, opts)
	if err != nil {
		return err
	}

	d.SetId(uuid)

	poolsList := []map[string]interface{}{}
	for _, x := range pools.Body {
		pool := make(map[string]interface{})
		pool["contract_number"] = x.ContractNumber
		pool["end_date"] = x.EndDate
		pool["entitlements_available"] = x.EntitlementsAvailable
		pool["id"] = x.Id
		pool["service_level"] = x.ServiceLevel
		pool["sku"] = x.Sku
		pool["start_date"] = x.StartDate
		pool["subscription_name"] = x.SubscriptionName
		pool["subscription_number"] = x.SubscriptionNumber

		poolsList = append(poolsList, pool)

	}
	d.Set("pools", poolsList)

	return nil
}
