package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CloudAccessDataSource{}

func NewCloudAccessDataSource() datasource.DataSource {
	return &CloudAccessDataSource{}
}

// CloudAccessDataSource defines the data source implementation.
type CloudAccessDataSource struct {
	client *apiClient
}

// CloudAccessDataSourceModel describes the data source data model.
type CloudAccessDataSourceModel struct {
	EnabledAccounts types.List `tfsdk:"enabled_accounts"`
}

func (m CloudAccessDataSourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled_accounts": types.ListType{ElemType: types.ObjectType{AttrTypes: EnabledAccountsModel{}.AttributeTypes()}},
	}
}

type EnabledAccountsModel struct {
	Accounts  types.List   `tfsdk:"accounts"`
	Name      types.String `tfsdk:"name"`
	Products  types.List   `tfsdk:"products"`
	ShortName types.String `tfsdk:"short_name"`
}

func (m EnabledAccountsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"accounts":   types.ObjectType{AttrTypes: AccountsModel{}.AttributeTypes()},
		"name":       types.StringType,
		"products":   types.ObjectType{AttrTypes: ProductsModel{}.AttributeTypes()},
		"short_name": types.StringType,
	}
}

type AccountsModel struct {
	DateAdded       types.String `tfsdk:"date_added"`
	GoldImageStatus types.List   `tfsdk:"gold_image_status"`
	ID              types.String `tfsdk:"id"`
	Nickname        types.String `tfsdk:"nickname"`
	SourceID        types.String `tfsdk:"source_id"`
	Verified        types.Bool   `tfsdk:"verified"`
}

func (m AccountsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"date_added":        types.StringType,
		"gold_image_status": types.ListType{ElemType: types.ObjectType{AttrTypes: GoldImageStatusModel{}.AttributeTypes()}},
		"id":                types.StringType,
		"nickname":          types.StringType,
		"source_id":         types.StringType,
		"verified":          types.BoolType,
	}
}

type ProductsModel struct {
	EnabledQuantity types.Int64  `tfsdk:"enabled_quantity"`
	ImageGroups     types.List   `tfsdk:"image_groups"`
	Name            types.String `tfsdk:"name"`
	NextRenewal     types.String `tfsdk:"next_renewal"`
	SKU             types.String `tfsdk:"sku"`
	TotalQuantity   types.Int64  `tfsdk:"total_quantity"`
}

func (m ProductsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled_quantity": types.Int64Type,
		"image_groups":     types.ListType{ElemType: types.StringType},
		"name":             types.StringType,
		"next_renewal":     types.StringType,
		"sku":              types.StringType,
		"total_quantity":   types.Int64Type,
	}
}

func (d *CloudAccessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_access"
}

func (d *CloudAccessDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to look up information about cloud providers entitled to Red Hat Cloud Access.",

		Attributes: map[string]schema.Attribute{
			"enabled_accounts": schema.ListNestedAttribute{
				MarkdownDescription: "A list where each entry is a single cloud provider",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"accounts": schema.ListNestedAttribute{
							Description: "A list of cloud accounts that are enabled for cloud access in the cloud provider.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"date_added": schema.StringAttribute{
										Description: "The date the account was added to cloud access.",
										Computed:    true,
									},
									"gold_image_status": schema.ListNestedAttribute{
										Description: "The status of any requests for gold image access for a cloud account.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"description": schema.StringAttribute{
													Description: "The description of the gold image.",
													Computed:    true,
												},
												"name": schema.StringAttribute{
													Description: "The name of the gold image.",
													Computed:    true,
												},
												"status": schema.StringAttribute{
													Description: "The status of the gold image request.",
													Computed:    true,
												},
											},
										},
									},
									"id": schema.StringAttribute{
										Description: "The id of the cloud account.",
										Computed:    true,
									},
									"nickname": schema.StringAttribute{
										Description: "A nickname associated with the cloud account.",
										Computed:    true,
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
							},
						},
						"name": schema.StringAttribute{
							Description: "The name of the cloud provider.",
							Computed:    true,
						},
						"products": schema.ListNestedAttribute{
							Description: "A list of products that are entitled to the cloud provider.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"enabled_quantity": schema.Int64Attribute{
										Description: "The quantity of subscriptions allowed to be consumed by the cloud provider.",
										Computed:    true,
									},
									"image_groups": schema.ListAttribute{
										Description: "A list of images associated with the cloud provider. These are used when requesting access to gold images for a cloud account.",
										Computed:    true,
										ElementType: types.StringType,
									},
									"name": schema.StringAttribute{
										Description: "The name of the product.",
										Computed:    true,
									},
									"next_renewal": schema.StringAttribute{
										Description: "The renewal date of the subscription.",
										Computed:    true,
									},
									"sku": schema.StringAttribute{
										Description: "The SKU of the product.",
										Computed:    true,
									},
									"total_quantity": schema.Int64Attribute{
										Description: "The total number of subscriptions of the product available.",
										Computed:    true,
									},
								},
							},
						},
						"short_name": schema.StringAttribute{
							Description: "An abreviation of the cloud provider name. Used when adding or removing accounts.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudAccessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*apiClient)
}

func (d *CloudAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudAccessDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := d.client.Client
	auth := d.client.Auth

	cap, _, err := client.CloudaccessAPI.ListEnabledCloudAccessProviders(auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to list enabled cloud access providers", err.Error())
		return
	}

	cloudProviders := []EnabledAccountsModel{}

	for _, x := range cap.GetBody() {
		// cloudProvider := make(map[string]interface{})
		cloudProvider := EnabledAccountsModel{
			Name:      types.StringValue(x.GetName()),
			ShortName: types.StringValue(x.GetShortName()),
		}

		// accounts := make([]map[string]interface{}, 0)
		accounts := []AccountsModel{}
		for _, y := range x.GetAccounts() {
			// account := make(map[string]interface{})
			account := AccountsModel{
				DateAdded: types.StringValue(y.GetDateAdded()),
				ID:        types.StringValue(y.GetId()),
				Nickname:  types.StringValue(y.GetNickname()),
				SourceID:  types.StringValue(y.GetSourceId()),
				Verified:  types.BoolValue(y.GetVerified()),
			}

			// goldImages := make([]map[string]interface{}, 0)
			goldImages := []GoldImageStatusModel{}
			if y.GoldImageStatus != nil {
				for _, z := range y.GetGoldImageStatus() {
					// goldImage := make(map[string]interface{})
					goldImage := GoldImageStatusModel{
						Description: types.StringValue(z.GetDescription()),
						Name:        types.StringValue(z.GetName()),
						Status:      types.StringValue(z.GetStatus()),
					}
					goldImages = append(goldImages, goldImage)
				}
			}
			goldImageStatus, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: GoldImageStatusModel{}.AttributeTypes()}, goldImages)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}

			account.GoldImageStatus = goldImageStatus
			accounts = append(accounts, account)
		}
		accountsList, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AccountsModel{}.AttributeTypes()}, accounts)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		cloudProvider.Accounts = accountsList

		products := []ProductsModel{}
		for _, y := range x.GetProducts() {
			// product := make(map[string]interface{})
			product := ProductsModel{
				EnabledQuantity: types.Int64Value(int64(y.GetEnabledQuantity())),
				Name:            types.StringValue(y.GetName()),
				SKU:             types.StringValue(y.GetSku()),
				TotalQuantity:   types.Int64Value(int64(y.GetTotalQuantity())),
				NextRenewal:     types.StringValue(y.GetNextRenewal()),
			}

			imageGroups, diag := types.ListValueFrom(ctx, types.StringType, y.GetImageGroups())
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			product.ImageGroups = imageGroups
			products = append(products, product)
		}
		productsList, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ProductsModel{}.AttributeTypes()}, products)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		cloudProvider.Products = productsList

		cloudProviders = append(cloudProviders, cloudProvider)
	}

	enabledCloudProviders, diag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: EnabledAccountsModel{}.AttributeTypes()}, cloudProviders)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	data.EnabledAccounts = enabledCloudProviders

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
