package rhsm

import (
	"regexp"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/umich-vci/gorhsm"
)

var nameRegex, _ = regexp.Compile("^[a-zA-Z0-9\\_\\-\\.]{1,100}$")

func resourceAllocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAllocationCreate,
		Read:   resourceAllocationRead,
		//Update: resourceAllocationUpdate,
		Delete: resourceAllocationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(nameRegex, "name must be less than 100 characters and can use only numbers, letters, underscores, hyphens, and periods"),
			},
			"uuid": &schema.Schema{
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

func resourceAllocationRead(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	uuid := d.Id()

	optional.NewString("entitlements")

	opts := &gorhsm.ShowAllocationOpts{
		Include: optional.NewString("entitlements"),
	}

	alloc, resp, err := client.AllocationApi.ShowAllocation(auth, uuid, opts)
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(alloc.Body.Uuid)
	d.Set("name", alloc.Body.Name)
	d.Set("uuid", alloc.Body.Uuid)
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

func resourceAllocationCreate(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	alloc, _, err := client.AllocationApi.CreateSatellite(auth, name)
	if err != nil {
		return err
	}

	d.SetId(alloc.Body.Uuid)

	return resourceAllocationRead(d, meta)
}

func resourceAllocationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAllocationRead(d, meta)
}

func resourceAllocationDelete(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	uuid := d.Id()

	_, err = client.AllocationApi.RemoveAllocation(auth, uuid, true)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
