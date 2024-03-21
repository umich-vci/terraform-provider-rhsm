package provider

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/umich-vci/gorhsm"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CloudAccessAccountResource{}
var _ resource.ResourceWithImportState = &CloudAccessAccountResource{}

func NewCloudAccessAccountResource() resource.Resource {
	return &CloudAccessAccountResource{}
}

// CloudAccessAccountResource defines the resource implementation.
type CloudAccessAccountResource struct {
	client *apiClient
}

// CloudAccessAccountResourceModel describes the resource data model.
type CloudAccessAccountResourceModel struct {
	ID                types.String `tfsdk:"id"`
	AccountID         types.String `tfsdk:"account_id"`
	ProviderShortName types.String `tfsdk:"provider_short_name"`
	GoldImages        types.Set    `tfsdk:"gold_images"`
	Nickname          types.String `tfsdk:"nickname"`
	DateAdded         types.String `tfsdk:"date_added"`
	GoldImageStatus   types.Set    `tfsdk:"gold_image_status"`
	SourceID          types.String `tfsdk:"source_id"`
	Verified          types.Bool   `tfsdk:"verified"`
}

type GoldImageStatusModel struct {
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
}

func (m GoldImageStatusModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description": types.StringType,
		"name":        types.StringType,
		"status":      types.StringType,
	}
}

func (r *CloudAccessAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_access_account"
}

func (r *CloudAccessAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource to manage entitlement for Red Hat Cloud Access for an account in a supported cloud provider.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the cloud account in the format `provider_short_name:account_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Description: "The ID of a cloud account that you would like to request Red Hat Cloud Access for. For GCE this should be a Google Group.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.NoneOf("")},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_short_name": schema.StringAttribute{
				Description: "The short name of the cloud provider that the `account_id` is in. This must be one of \"AWS\", \"GCE\", or \"MSAZ\".  Other cloud providers are supported but have not been tested so they are not in the list of valid options.",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf(cloudAccessAccountProviders...)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gold_images": schema.SetAttribute{
				Description: "A list of gold images to request access to for the account. Images available to a cloud provider can be found with the `rhsm_cloud_access` data source. Once you request access to a gold image, it is not possible to disable access via the API.",
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				ElementType: types.StringType,
			},
			"nickname": schema.StringAttribute{
				Description: "A nickname to help describe the account.",
				Optional:    true,
			},
			"date_added": schema.StringAttribute{
				Description: "The date the cloud account was added to Red Hat Cloud Access.",
				Computed:    true,
			},
			"gold_image_status": schema.SetAttribute{
				Description: "The status of any requests for gold image access for the cloud account.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"description": types.StringType,
						"name":        types.StringType,
						"status":      types.StringType,
					},
				},
			},
			"source_id": schema.StringAttribute{
				Description: "Source ID of linked account. Only for accounts created via Sources on cloud.redhat.com.",
				Computed:    true,
			},
			"verified": schema.BoolAttribute{
				Description: "Is the cloud provider account verified for RHSM Auto Registration?",
				Computed:    true,
			},
		},
	}
}

func (r *CloudAccessAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*apiClient)
}

func (r *CloudAccessAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudAccessAccountResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.Client
	auth := r.client.Auth

	data.ID = types.StringValue(fmt.Sprintf("%s:%s", data.ProviderShortName.ValueString(), data.AccountID.ValueString()))

	account := &gorhsm.AddProviderAccount{
		Id:       data.AccountID.ValueStringPointer(),
		Nickname: data.Nickname.ValueStringPointer(),
	}
	accountList := []gorhsm.AddProviderAccount{*account}

	apa, err := client.CloudaccessAPI.AddProviderAccounts(auth, data.ProviderShortName.ValueString()).Account(accountList).Execute()
	if apa != nil {
		defer apa.Body.Close()
	}
	if err != nil {
		apaBody, e := io.ReadAll(apa.Body)
		if e != nil && apaBody != nil {
			resp.Diagnostics.AddError("Failed to create Cloud Access Account", err.Error())

		} else {
			resp.Diagnostics.AddError("Failed to create Cloud Access Account", string(apaBody))
		}
		return
	}

	var goldImages []string
	data.GoldImages.ElementsAs(ctx, &goldImages, false)

	// do not enable gold images if none are specified
	if len(goldImages) > 0 {
		gi := &gorhsm.EnableGoldImagesRequest{
			Accounts: []string{data.AccountID.ValueString()},
			Images:   goldImages,
		}

		egi, err := client.CloudaccessAPI.EnableGoldImages(auth, data.ProviderShortName.ValueString()).GoldImages(*gi).Execute()
		if egi != nil {
			defer egi.Body.Close()
		}
		if err != nil {
			egiBody, e := io.ReadAll(egi.Body)
			if e != nil && egiBody != nil {
				resp.Diagnostics.AddError("Failed to enable gold images", err.Error())
			} else {
				resp.Diagnostics.AddError("Failed to enable gold images", string(egiBody))
			}
			return
		}
	}
	caps, capsRaw, err := client.CloudaccessAPI.ListEnabledCloudAccessProviders(auth).Execute()
	if capsRaw != nil {
		defer capsRaw.Body.Close()
	}
	if err != nil {
		capsBody, e := io.ReadAll(capsRaw.Body)
		if e != nil && capsBody != nil {
			resp.Diagnostics.AddError("Failed to list enabled cloud access providers", err.Error())
		} else {
			resp.Diagnostics.AddError("Failed to list enabled cloud access providers", string(capsBody))
		}
		return
	}

	caa, diag := flattenCloudAccessAccount(ctx, caps, data.ProviderShortName.ValueString(), data.AccountID.ValueString())
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// no matching account was found which should not happen since we just created it
	if caa == nil {
		resp.Diagnostics.AddError("Failed to find created Cloud Access Account", "No matching account was found")
		return
	}

	data.AccountID = caa.AccountID
	data.ProviderShortName = caa.ProviderShortName
	data.Nickname = caa.Nickname
	data.DateAdded = caa.DateAdded
	data.SourceID = caa.SourceID
	data.Verified = caa.Verified
	data.GoldImages = caa.GoldImages
	data.GoldImageStatus = caa.GoldImageStatus

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudAccessAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudAccessAccountResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.Client
	auth := r.client.Auth

	shortName, accountID, err := resourceCloudAccessAccountSplitID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse Cloud Access Account resource ID", err.Error())
		return
	}

	caps, _, err := client.CloudaccessAPI.ListEnabledCloudAccessProviders(auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to list enabled cloud access providers", err.Error())
		return
	}

	caa, diag := flattenCloudAccessAccount(ctx, caps, shortName, accountID)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// no matching account was found and no error was returned
	if caa == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.AccountID = caa.AccountID
	data.ProviderShortName = caa.ProviderShortName
	data.Nickname = caa.Nickname
	data.DateAdded = caa.DateAdded
	data.SourceID = caa.SourceID
	data.Verified = caa.Verified
	data.GoldImages = caa.GoldImages
	data.GoldImageStatus = caa.GoldImageStatus

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudAccessAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state CloudAccessAccountResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.Client
	auth := r.client.Auth

	shortName, accountID, err := resourceCloudAccessAccountSplitID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse Cloud Access Account resource ID", err.Error())
		return
	}

	if !data.Nickname.Equal(state.Nickname) {
		account := &gorhsm.UpdateProviderAccountRequest{Nickname: data.Nickname.ValueString()}
		_, err := client.CloudaccessAPI.UpdateProviderAccount(auth, shortName, accountID).Account(*account).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Failed to update Cloud Access Account nickname", err.Error())
			return
		}
	}

	if !data.GoldImages.Equal(state.GoldImages) {
		goldImages := []string{}
		resp.Diagnostics.Append(data.GoldImages.ElementsAs(ctx, goldImages, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		gi := &gorhsm.EnableGoldImagesRequest{
			Accounts: []string{accountID},
			Images:   goldImages,
		}

		_, err := client.CloudaccessAPI.EnableGoldImages(auth, shortName).GoldImages(*gi).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Failed to enable gold images", err.Error())
			return
		}
	}

	caps, _, err := client.CloudaccessAPI.ListEnabledCloudAccessProviders(auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to list enabled cloud access providers", err.Error())
		return
	}

	caa, diag := flattenCloudAccessAccount(ctx, caps, data.ProviderShortName.ValueString(), data.AccountID.ValueString())
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// no matching account was found which should not happen since we just updated it
	if caa == nil {
		resp.Diagnostics.AddError("Failed to find updated Cloud Access Account", "No matching account was found")
		return
	}

	data.AccountID = caa.AccountID
	data.ProviderShortName = caa.ProviderShortName
	data.Nickname = caa.Nickname
	data.DateAdded = caa.DateAdded
	data.SourceID = caa.SourceID
	data.Verified = caa.Verified
	data.GoldImages = caa.GoldImages
	data.GoldImageStatus = caa.GoldImageStatus

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudAccessAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudAccessAccountResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.Client
	auth := r.client.Auth

	shortName, accountID, err := resourceCloudAccessAccountSplitID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse Cloud Access Account resource ID", err.Error())
		return
	}

	remove := &gorhsm.RemoveProviderAccountRequest{
		Id: accountID,
	}

	_, err = client.CloudaccessAPI.RemoveProviderAccount(auth, shortName).Account(*remove).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to remove Cloud Access Account", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CloudAccessAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type CloudAccessAccountModel struct {
	ID                types.String
	AccountID         types.String
	ProviderShortName types.String
	GoldImages        types.Set
	Nickname          types.String
	DateAdded         types.String
	GoldImageStatus   types.Set
	SourceID          types.String
	Verified          types.Bool
}

func flattenCloudAccessAccount(ctx context.Context, caps *gorhsm.ListEnabledCloudAccessProviders200Response, shortName string, accountID string) (*CloudAccessAccountModel, diag.Diagnostics) {
	var d diag.Diagnostics

	if caps == nil {
		return nil, d
	}

	caa := new(CloudAccessAccountModel)

	for _, x := range caps.GetBody() {
		if x.GetShortName() == shortName {
			for _, y := range x.GetAccounts() {
				if y.Id == accountID {
					goldImageStatus := []types.Object{}
					goldImages := []attr.Value{}
					for _, z := range y.GetGoldImageStatus() {
						// goldImage := make(map[string]attr.Value)
						// goldImage["description"] = types.StringValue(z.GetDescription())
						// goldImage["name"] = types.StringValue(z.GetName())
						// goldImage["status"] = types.StringValue(z.GetStatus())
						goldImage := GoldImageStatusModel{
							Description: types.StringValue(z.GetDescription()),
							Name:        types.StringValue(z.GetDescription()),
							Status:      types.StringValue(*z.Status),
						}
						goldImageObject, diag := types.ObjectValueFrom(ctx, goldImage.AttributeTypes(), goldImage)
						if diag.HasError() {
							d.Append(diag...)
							return nil, d
						}
						goldImageStatus = append(goldImageStatus, goldImageObject)
						goldImages = append(goldImages, types.StringValue(z.GetName()))
					}

					goldImageStatusSet, diag := types.SetValueFrom(ctx,
						types.ObjectType{AttrTypes: GoldImageStatusModel{}.AttributeTypes()}, goldImageStatus)
					if diag.HasError() {
						d.Append(diag...)
						return nil, d
					}

					goldImagesSet, diag := types.SetValue(types.StringType, goldImages)
					if diag.HasError() {
						d.Append(diag...)
						return nil, d
					}

					caa = &CloudAccessAccountModel{
						ID:                types.StringValue(fmt.Sprintf("%s:%s", shortName, accountID)),
						AccountID:         types.StringValue(accountID),
						ProviderShortName: types.StringValue(shortName),
						Nickname:          types.StringValue(y.GetNickname()),
						DateAdded:         types.StringValue(y.GetDateAdded()),
						SourceID:          types.StringValue(y.GetSourceId()),
						Verified:          types.BoolValue(y.GetVerified()),
						GoldImages:        goldImagesSet,
						GoldImageStatus:   goldImageStatusSet,
					}
					break
				}
			}
			break
		}
	}
	return caa, d
}

func resourceCloudAccessAccountSplitID(id string) (shortName string, accountID string, err error) {
	splitID := strings.SplitN(id, ":", 2)

	if len(splitID) != 2 {
		return "", "", fmt.Errorf("the Cloud Access Account ID %s could not be split correctly", id)
	}

	name := splitID[0]
	acctID := splitID[1]

	validProvider := false
	for _, x := range cloudAccessAccountProviders {
		if name == x {
			validProvider = true
		}
	}

	if !validProvider {
		return "", "", fmt.Errorf("invalid Cloud Access Account provider %s specified in ID %s", name, id)
	}

	return name, acctID, nil
}

var cloudAccessAccountProviders = []string{
	"AWS",
	"GCE",
	"MSAZ",
}
