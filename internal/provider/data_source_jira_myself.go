package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraMyselfDataSource struct {
		p atlassianProvider
	}

	jiraMyselfDataSourceModel struct {
		ID               types.String                      `tfsdk:"id"`
		Self             types.String                      `tfsdk:"self"`
		AccountId        types.String                      `tfsdk:"account_id"`
		AccountType      types.String                      `tfsdk:"account_type"`
		EmailAddress     types.String                      `tfsdk:"email_address"`
		AvatarUrls       *jiraMyselfAvatarUrlsModel        `tfsdk:"avatar_urls"`
		DisplayName      types.String                      `tfsdk:"display_name"`
		Active           types.Bool                        `tfsdk:"active"`
		TimeZone         types.String                      `tfsdk:"timezone"`
		Locale           types.String                      `tfsdk:"locale"`
		Groups           []jiraMyselfGroupsModel           `tfsdk:"groups"`
		ApplicationRoles []jiraMyselfApplicationRolesModel `tfsdk:"application_roles"`
	}

	jiraMyselfAvatarUrlsModel struct {
		One6X16   types.String `tfsdk:"p16x16"`
		Two4X24   types.String `tfsdk:"p24x24"`
		Three2X32 types.String `tfsdk:"p32x32"`
		Four8X48  types.String `tfsdk:"p48x48"`
	}

	jiraMyselfGroupsModel struct {
		Name types.String `tfsdk:"name"`
		Self types.String `tfsdk:"self"`
	}

	jiraMyselfApplicationRolesModel struct {
		Key                  types.String `tfsdk:"key"`
		Groups               types.List   `tfsdk:"groups"`
		Name                 types.String `tfsdk:"name"`
		DefaultGroups        types.List   `tfsdk:"default_groups"`
		SelectedByDefault    types.Bool   `tfsdk:"select_by_default"`
		NumberOfSeats        types.Int64  `tfsdk:"number_of_seats"`
		RemainingSeats       types.Int64  `tfsdk:"remaining_seats"`
		UserCount            types.Int64  `tfsdk:"user_count"`
		UserCountDescription types.String `tfsdk:"user_count_description"`
		HasUnlimitedSeats    types.Bool   `tfsdk:"has_unlimited_seats"`
		Platform             types.Bool   `tfsdk:"platform"`
	}
)

var (
	_ datasource.DataSource = (*jiraMyselfDataSource)(nil)
)

func NewJiraMyselfDataSource() datasource.DataSource {
	return &jiraMyselfDataSource{}
}

func (*jiraMyselfDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_myself"
}

func (*jiraMyselfDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Myself Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of myself. Defaults to value of `account_id`.",
				Computed:            true,
				Type:                types.StringType,
			},
			"self": {
				MarkdownDescription: "The URL of the user.",
				Computed:            true,
				Type:                types.StringType,
			},
			"account_id": {
				MarkdownDescription: "The account ID of the user, which uniquely identifies the user across all Atlassian products.",
				Computed:            true,
				Type:                types.StringType,
			},
			"account_type": {
				MarkdownDescription: "The user account type. Can take the following values: `atlassian`, `app`, `customer`.",
				Computed:            true,
				Type:                types.StringType,
			},
			"email_address": {
				MarkdownDescription: "The email address of the user. Depending on the user’s privacy setting, this may be returned as null.",
				Computed:            true,
				Type:                types.StringType,
			},
			"avatar_urls": {
				MarkdownDescription: "The avatars of the user.",
				Computed:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"p16x16": {
						MarkdownDescription: "The URL of the item's 16x16 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
					},
					"p24x24": {
						MarkdownDescription: "The URL of the item's 24x24 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
					},
					"p32x32": {
						MarkdownDescription: "The URL of the item's 32x32 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
					},
					"p48x48": {
						MarkdownDescription: "The URL of the item's 48x48 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
					},
				}),
			},
			"display_name": {
				MarkdownDescription: "The display name of the user. Depending on the user’s privacy setting, this may return an alternative value.",
				Computed:            true,
				Type:                types.StringType,
			},
			"active": {
				MarkdownDescription: "Whether the user is active.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"timezone": {
				MarkdownDescription: "The time zone specified in the user's profile. Depending on the user’s privacy setting, this may be returned as null.",
				Computed:            true,
				Type:                types.StringType,
			},
			"locale": {
				MarkdownDescription: "The locale of the user. Depending on the user’s privacy setting, this may be returned as null.",
				Computed:            true,
				Type:                types.StringType,
			},
			"groups": {
				MarkdownDescription: "The groups that the user belongs to.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(
					map[string]tfsdk.Attribute{
						"name": {
							MarkdownDescription: "The name of the group.",
							Computed:            true,
							Type:                types.StringType,
						},
						"self": {
							MarkdownDescription: "The URL for the group details.",
							Computed:            true,
							Type:                types.StringType,
						},
					},
				),
			},
			"application_roles": {
				MarkdownDescription: "The application roles the user is assigned to.",
				Computed:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"key": {
						MarkdownDescription: "The key of the application role.",
						Computed:            true,
						Type:                types.StringType,
					},
					"groups": {
						MarkdownDescription: "The groups associated with the application role.",
						Computed:            true,
						Type:                types.ListType{ElemType: types.StringType},
					},
					"name": {
						MarkdownDescription: "The display name of the application role.",
						Computed:            true,
						Type:                types.StringType,
					},
					"default_groups": {
						MarkdownDescription: "The groups that are granted default access for this application role.",
						Computed:            true,
						Type:                types.ListType{ElemType: types.StringType},
					},
					"select_by_default": {
						MarkdownDescription: "Determines whether this application role should be selected by default on user creation.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"number_of_seats": {
						MarkdownDescription: "The maximum count of users on your license.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"remaining_seats": {
						MarkdownDescription: "The count of users remaining on your license.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"user_count": {
						MarkdownDescription: "The number of users counting against your license.",
						Computed:            true,
						Type:                types.Int64Type,
					},
					"user_count_description": {
						MarkdownDescription: "The type of users being counted against your license.",
						Computed:            true,
						Type:                types.StringType,
					},
					"has_unlimited_seats": {
						MarkdownDescription: "Whether unlimited user licenses are available.",
						Computed:            true,
						Type:                types.BoolType,
					},
					"platform": {
						MarkdownDescription: "Indicates if the application role belongs to Jira platform (jira-core).",
						Computed:            true,
						Type:                types.BoolType,
					},
				}),
			},
		},
	}, nil
}

func (d *jiraMyselfDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jira.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jira.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.p.jira = client
}

func (d *jiraMyselfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading myself data source")

	myself, res, err := d.p.jira.MySelf.Details(ctx, []string{"groups", "applicationRoles"})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get myself, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved myself from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v, groups:%+v, applicationRoles:%+v", myself, *myself.Groups, *myself.ApplicationRoles),
	})

	newState := jiraMyselfDataSourceModel{
		ID:           types.StringValue(myself.AccountID),
		Self:         types.StringValue(myself.Self),
		AccountId:    types.StringValue(myself.AccountID),
		AccountType:  types.StringValue(myself.AccountType),
		EmailAddress: types.StringValue(myself.EmailAddress),
		AvatarUrls: &jiraMyselfAvatarUrlsModel{
			One6X16:   types.StringValue(myself.AvatarUrls.One6X16),
			Two4X24:   types.StringValue(myself.AvatarUrls.Two4X24),
			Three2X32: types.StringValue(myself.AvatarUrls.Three2X32),
			Four8X48:  types.StringValue(myself.AvatarUrls.Four8X48),
		},
		DisplayName:      types.StringValue(myself.DisplayName),
		Active:           types.BoolValue(myself.Active),
		TimeZone:         types.StringValue(myself.TimeZone),
		Locale:           types.StringValue(myself.Locale),
		Groups:           []jiraMyselfGroupsModel{},
		ApplicationRoles: []jiraMyselfApplicationRolesModel{},
	}

	// Get groups
	var groups []jiraMyselfGroupsModel
	for _, v := range myself.Groups.Items {
		g := jiraMyselfGroupsModel{
			Name: types.StringValue(v.Name),
			Self: types.StringValue(v.Self),
		}

		groups = append(groups, g)
	}
	newState.Groups = groups

	// Get applicationroles
	var roles []jiraMyselfApplicationRolesModel
	for _, v := range myself.ApplicationRoles.Items {
		r := jiraMyselfApplicationRolesModel{
			Key:                  types.StringValue(v.Key),
			Name:                 types.StringValue(v.Name),
			Groups:               types.ListNull(types.StringType),
			DefaultGroups:        types.ListNull(types.StringType),
			SelectedByDefault:    types.BoolValue(v.SelectedByDefault),
			NumberOfSeats:        types.Int64Value(int64(v.NumberOfSeats)),
			RemainingSeats:       types.Int64Value(int64(v.RemainingSeats)),
			UserCount:            types.Int64Value(int64(v.UserCount)),
			UserCountDescription: types.StringValue(v.UserCountDescription),
			HasUnlimitedSeats:    types.BoolValue(v.HasUnlimitedSeats),
			Platform:             types.BoolValue(v.Platform),
		}
		// Get groups
		r.Groups, _ = types.ListValueFrom(ctx, types.StringType, v.Groups)
		// Get default groups
		r.DefaultGroups, _ = types.ListValueFrom(ctx, types.StringType, v.DefaultGroups)

		roles = append(roles, r)
	}
	newState.ApplicationRoles = roles

	tflog.Debug(ctx, "Storing myself into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
