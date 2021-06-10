package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gorhsm"
)

func resourceCloudAccessAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage entitlement for Red Hat Cloud Access for an account in a supported cloud provider.",

		CreateContext: resourceCloudAccessAccountCreate,
		ReadContext:   resourceCloudAccessAccountRead,
		UpdateContext: resourceCloudAccessAccountUpdate,
		DeleteContext: resourceCloudAccessAccountDelete,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "The ID of a cloud account that you would like to request Red Hat Cloud Access for.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
			},
			"provider_short_name": {
				Description:  "The short name of the cloud provider that the `account_id` is in. This must be one of \"AWS\", \"GCE\", or \"MSAZ\".  Other cloud providers are supported but have not been tested so they are not in the list of valid options.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "GCE", "MSAZ"}, false),
			},
			"gold_images": {
				Description: "A list of gold images to request access to for the account. Images available to a cloud provider can be found with the `rhsm_cloud_access` data source.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"nickname": {
				Description: "A nickname to help describe the account.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"date_added": {
				Description: "The date the cloud account was added to Red Hat Cloud Access.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"gold_image_status": {
				Description: "The status of any requests for gold image access for the cloud account.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Description: "The description of the gold image.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the gold image.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status": {
							Description: "The status of the gold image request.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"source_id": {
				Description: "Source ID of linked account. Only for accounts created via Sources on cloud.redhat.com.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"verified": {
				Description: "Is the cloud provider account verified for RHSM Auto Registration?",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourceCloudAccessAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	id := d.Get("account_id").(string)
	shortName := d.Get("provider_short_name").(string)
	foundAccount := false

	cap, _, err := client.CloudaccessApi.ListEnabledCloudAccessProviders(auth).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, x := range cap.GetBody() {
		if x.GetShortName() == shortName {
			for _, y := range x.GetAccounts() {
				if y.Id == id {
					foundAccount = true
					d.Set("nickname", y.GetNickname())
					d.Set("date_added", y.GetDateAdded())
					d.Set("source_id", y.GetSourceId())
					d.Set("verified", y.GetVerified())

					goldImageStatus := make([]map[string]interface{}, 0)
					for _, z := range y.GetGoldImageStatus() {
						goldImage := make(map[string]interface{})
						goldImage["description"] = z.GetDescription()
						goldImage["name"] = z.GetName()
						goldImage["status"] = z.GetStatus()
						goldImageStatus = append(goldImageStatus, goldImage)
					}
					d.Set("gold_image_status", goldImageStatus)
					break
				}
			}
		}
	}

	if !foundAccount {
		d.SetId("")
	}

	return nil
}

func resourceCloudAccessAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	id := d.Get("account_id").(string)
	shortName := d.Get("provider_short_name").(string)
	nickname := d.Get("nickname").(string)

	account := &gorhsm.AddProviderAccount{
		Id:       &id,
		Nickname: &nickname,
	}
	accountList := []gorhsm.AddProviderAccount{*account}

	_, err := client.CloudaccessApi.AddProviderAccounts(auth, shortName).Account(accountList).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	if g, ok := d.GetOk("gold_images"); ok {
		rawGoldImages := g.(*schema.Set).List()
		goldimages := []string{}
		for x := range rawGoldImages {
			goldimages = append(goldimages, rawGoldImages[x].(string))
		}
		gi := &gorhsm.InlineObject5{
			Accounts: []string{id},
			Images:   goldimages,
		}

		_, err = client.CloudaccessApi.EnableGoldImages(auth, shortName).GoldImages(*gi).Execute()
		if err != nil {
			d.Set("gold_images", []string{})
			return diag.FromErr(err)
		}
	}

	return resourceCloudAccessAccountRead(ctx, d, meta)
}

func resourceCloudAccessAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	id := d.Id()
	shortName := d.Get("provider_short_name").(string)
	accountID := d.Get("account_id").(string)

	if d.HasChange("nickname") {
		account := &gorhsm.InlineObject3{Nickname: d.Get("nickname").(string)}
		_, err := client.CloudaccessApi.UpdateProviderAccount(auth, shortName, accountID).Account(*account).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("gold_images") {
		if g, ok := d.GetOk("gold_images"); ok {
			rawGoldImages := g.(*schema.Set).List()
			goldimages := []string{}
			for x := range rawGoldImages {
				goldimages = append(goldimages, rawGoldImages[x].(string))
			}
			gi := &gorhsm.InlineObject5{
				Accounts: []string{id},
				Images:   goldimages,
			}

			_, err := client.CloudaccessApi.EnableGoldImages(auth, shortName).GoldImages(*gi).Execute()
			if err != nil {
				d.Set("gold_images", []string{})
				return diag.FromErr(err)
			}
		}
	}

	return resourceCloudAccessAccountRead(ctx, d, meta)
}

func resourceCloudAccessAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client
	auth := meta.(*apiClient).Auth

	id := d.Id()
	shortName := d.Get("provider_short_name").(string)

	remove := &gorhsm.InlineObject2{
		Id: id,
	}

	_, err := client.CloudaccessApi.RemoveProviderAccount(auth, shortName).Account(*remove).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
