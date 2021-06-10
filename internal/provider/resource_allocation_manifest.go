package provider

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gorhsm"
)

func resourceAllocationManifest() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to create a manifest from a RHSM subscription allocation that can be uploaded to a Red Hat Satellite server.",

		CreateContext: resourceAllocationManifestCreate,
		ReadContext:   resourceAllocationManifestRead,
		DeleteContext: resourceAllocationManifestDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocation_uuid": {
				Description:  "The UUID of the subscription allocation to create the manifest from.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"last_modified": {
				Description: "The date and time the subscription allocation was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"manifest_last_modified": {
				Description: "The date and time the subscription allocation was last modified when the manifest was last generated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"manifest": {
				Description: "The manifest as downloaded from the RHSM portal.  This is a zip file which has been base64 encoded to a string.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceAllocationManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	d.Set("last_modified", alloc.Body.GetLastModified())

	return nil
}

func resourceAllocationManifestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	allocationUUID := d.Get("allocation_uuid").(string)

	alloc, _, err := client.AllocationApi.ShowAllocation(auth, allocationUUID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("manifest_last_modified", alloc.Body.GetLastModified())

	exportJob, _, err := client.AllocationApi.ExportAllocation(auth, allocationUUID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	exportJobID := exportJob.Body.GetExportJobID()

	var manifestURL string
	for {
		time.Sleep(5 * time.Second)
		status, resp, err := client.AllocationApi.ExportJobAllocation(auth, allocationUUID, exportJobID).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode == 200 {
			manifestURL = status.Body.GetHref()
			break
		}
	}

	mclient := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, manifestURL, nil)
	req.Header.Add("Authorization", "Bearer "+auth.Value(gorhsm.ContextAPIKeys).(gorhsm.APIKey).Key)
	resp, err := mclient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	manifest := base64.StdEncoding.EncodeToString(body)
	d.Set("manifest", manifest)
	d.SetId(allocationUUID)

	return resourceAllocationManifestRead(ctx, d, meta)
}

func resourceAllocationManifestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
