package rhsm

import (
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gorhsm"
)

func resourceCloudAccessAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudAccessAccountCreate,
		Read:   resourceCloudAccessAccountRead,
		Update: resourceCloudAccessAccountUpdate,
		Delete: resourceCloudAccessAccountDelete,
		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"provider_short_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "GCE", "MSAZ"}, false),
			},
			"gold_images": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"nickname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"date_added": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"gold_image_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudAccessAccountRead(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	id := d.Get("account_id").(string)
	shortName := d.Get("provider_short_name").(string)
	foundAccount := false

	cap, _, err := client.CloudaccessApi.ListEnabledCloudAccessProviders(auth)
	if err != nil {
		return err
	}

	for _, x := range cap.Body {
		if x.ShortName == shortName {
			for _, y := range x.Accounts {
				if y.Id == id {
					foundAccount = true
					d.Set("nickname", y.Nickname)
					d.Set("date_added", y.DateAdded)
					d.Set("gold_image_status", y.GoldImageStatus)
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

func resourceCloudAccessAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	id := d.Get("account_id").(string)
	shortName := d.Get("provider_short_name").(string)
	nickname := d.Get("nickname").(string)

	account := &gorhsm.AddProviderAccount{
		Id:       id,
		Nickname: nickname,
	}
	accountList := []gorhsm.AddProviderAccount{*account}
	accountOpts := &gorhsm.AddProviderAccountsOpts{
		Account: optional.NewInterface(accountList),
	}

	_, err = client.CloudaccessApi.AddProviderAccounts(auth, shortName, accountOpts)
	if err != nil {
		return err
	}

	d.SetId(id)

	if g, ok := d.GetOk("gold_images"); ok {
		rawGoldImages := g.(*schema.Set).List()
		goldimages := []string{}
		for x := range rawGoldImages {
			goldimages = append(goldimages, rawGoldImages[x].(string))
		}
		gi := &gorhsm.InlineObject2{
			Accounts: []string{id},
			Images:   goldimages,
		}
		goldopts := &gorhsm.EnableGoldImagesOpts{
			GoldImages: optional.NewInterface(*gi),
		}

		_, err = client.CloudaccessApi.EnableGoldImages(auth, shortName, goldopts)
		if err != nil {
			d.Set("gold_images", []string{})
			return err
		}
	}

	return nil
}

func resourceCloudAccessAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	id := d.Id()
	shortName := d.Get("provider_short_name").(string)

	update := &gorhsm.InlineObject{
		Id: id,
	}

	idOrNameChange := false

	if d.HasChange("account_id") {
		idOrNameChange = true
		update.NewID = d.Get("account_id").(string)
	}

	if d.HasChange("nickname") {
		idOrNameChange = true
		update.NewNickname = d.Get("nickname").(string)
	}

	updateOpts := &gorhsm.UpdateProviderAccountOpts{
		Account: optional.NewInterface(*update),
	}

	if idOrNameChange {
		_, err = client.CloudaccessApi.UpdateProviderAccount(auth, shortName, updateOpts)
		if err != nil {
			return err
		}

		if d.HasChange("account_id") {
			d.SetId(d.Get("account_id").(string))
		}
	}

	if d.HasChange("gold_images") {
		if g, ok := d.GetOk("gold_images"); ok {
			rawGoldImages := g.(*schema.Set).List()
			goldimages := []string{}
			for x := range rawGoldImages {
				goldimages = append(goldimages, rawGoldImages[x].(string))
			}
			gi := &gorhsm.InlineObject2{
				Accounts: []string{id},
				Images:   goldimages,
			}
			goldopts := &gorhsm.EnableGoldImagesOpts{
				GoldImages: optional.NewInterface(*gi),
			}

			_, err = client.CloudaccessApi.EnableGoldImages(auth, shortName, goldopts)
			if err != nil {
				d.Set("gold_images", []string{})
				return err
			}
		}
	}

	return nil
}

func resourceCloudAccessAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client, auth, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	id := d.Id()
	shortName := d.Get("provider_short_name").(string)

	remove := &gorhsm.InlineObject1{
		Id: id,
	}
	removeOpts := &gorhsm.RemoveProviderAccountOpts{
		Account: optional.NewInterface(remove),
	}

	_, err = client.CloudaccessApi.RemoveProviderAccount(auth, shortName, removeOpts)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
