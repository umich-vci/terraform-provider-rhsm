package provider

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gorhsm"
)

func resourceAllocationManifest() *schema.Resource {
	return &schema.Resource{
		Create: resourceAllocationManifestCreate,
		Read:   resourceAllocationManifestRead,
		//Update: resourceAllocationManifestUpdate,
		Delete: resourceAllocationManifestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allocation_uuid": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"last_modified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manifest_last_modified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manifest": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAllocationManifestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	uuid := d.Id()

	alloc, resp, err := client.AllocationApi.ShowAllocation(auth, uuid).Execute()
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("last_modified", alloc.Body.LastModified)

	return nil
}

func resourceAllocationManifestCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)

	alloc, _, err := client.AllocationApi.ShowAllocation(auth, allocationUUID).Execute()
	if err != nil {
		return err
	}

	d.Set("manifest_last_modified", alloc.Body.LastModified)

	exportJob, _, err := client.AllocationApi.ExportAllocation(auth, allocationUUID).Execute()
	if err != nil {
		return err
	}

	exportJobID := *exportJob.Body.ExportJobID

	var manifestURL string
	for {
		time.Sleep(5 * time.Second)
		status, resp, err := client.AllocationApi.ExportJobAllocation(auth, allocationUUID, exportJobID).Execute()
		if err != nil {
			return err
		}
		if resp.StatusCode == 200 {
			manifestURL = *status.Body.Href
			break
		}
	}

	mclient := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, manifestURL, nil)
	req.Header.Add("Authorization", "Bearer "+auth.Value(gorhsm.ContextAPIKeys).(gorhsm.APIKey).Key)
	resp, err := mclient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	manifest := base64.StdEncoding.EncodeToString(body)
	d.Set("manifest", manifest)
	d.SetId(allocationUUID)

	return resourceAllocationManifestRead(d, meta)
}

func resourceAllocationManifestDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
