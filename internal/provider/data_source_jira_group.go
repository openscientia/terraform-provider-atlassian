package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	common "github.com/openscientia/terraform-provider-atlassian/internal/provider/models"
)

type (
	jiraGroupDataSource struct {
		p atlassianProvider
	}

	jiraGroupDataSourceModel struct {
		ID      types.String `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		GroupID types.String `tfsdk:"group_id"`
		Self    types.String `tfsdk:"self"`
		Users   types.Set    `tfsdk:"users"`
	}
)

var (
	_ datasource.DataSource = (*jiraGroupDataSource)(nil)
)

func NewJiraGroupDataSource() datasource.DataSource {
	return &jiraGroupDataSource{}
}

func (*jiraGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_group"
}

func (*jiraGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Jira Group Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group. Defaults to `group_id`.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the group.",
				Required:            true,
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group, which uniquely identifies the group across all Atlassian products.",
				Computed:            true,
			},
			"self": schema.StringAttribute{
				MarkdownDescription: "The URL for these group details.",
				Computed:            true,
			},
			"users": schema.SetNestedAttribute{
				MarkdownDescription: "The list of users in the group.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"self": schema.StringAttribute{
							MarkdownDescription: "The URL of the user.",
							Computed:            true,
						},
						"account_id": schema.StringAttribute{
							MarkdownDescription: "The account ID of the user, which uniquely identifies the user across all Atlassian products.",
							Computed:            true,
						},
						"email_address": schema.StringAttribute{
							MarkdownDescription: "The email address of the user. Depending on the user’s privacy settings, this may be returned as null.",
							Computed:            true,
						},
						"avatar_urls": schema.SingleNestedAttribute{
							MarkdownDescription: "The avatars of the user.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"p16x16": schema.StringAttribute{
									MarkdownDescription: "The URL of the item's 16x16 pixel avatar.",
									Computed:            true,
								},
								"p24x24": schema.StringAttribute{
									MarkdownDescription: "The URL of the item's 24x24 pixel avatar.",
									Computed:            true,
								},
								"p32x32": schema.StringAttribute{
									MarkdownDescription: "The URL of the item's 32x32 pixel avatar.",
									Computed:            true,
								},
								"p48x48": schema.StringAttribute{
									MarkdownDescription: "The URL of the item's 48x48 pixel avatar.",
									Computed:            true,
								},
							},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the user. Depending on the user’s privacy settings, this may return an alternative value.",
							Computed:            true,
						},
						"active": schema.BoolAttribute{
							MarkdownDescription: "Whether the user is active.",
							Computed:            true,
						},
						"timezone": schema.StringAttribute{
							MarkdownDescription: "The time zone specified in the user's profile. Depending on the user’s privacy settings, this may be returned as null.",
							Computed:            true,
						},
						"account_type": schema.StringAttribute{
							MarkdownDescription: "The type of account represented by this user. This will be one of `atlassian` (normal users), `app` (application user) or `customer` (Jira Service Desk customer user)",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *jiraGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading group data source")

	var newState jiraGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded group config", map[string]interface{}{
		"readConfig": fmt.Sprintf("%+v", newState),
	})

	opts := &models.GroupBulkOptionsScheme{
		GroupNames: []string{newState.Name.ValueString()},
	}
	group, res, err := d.p.jira.Group.Bulk(ctx, opts, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get group, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved group from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v", group.Values[0]),
	})

	isLast := false
	startAt := 0
	maxResults := 100
	members := []*models.GroupUserDetailScheme{}
	for !isLast {
		groupMembers, res, err := d.p.jira.Group.Members(ctx, newState.Name.ValueString(), true, startAt, maxResults)
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get group members, got error: %s\n%s", err, resBody))
			return
		}
		startAt += maxResults
		isLast = groupMembers.IsLast
		members = append(members, groupMembers.Values...)
	}
	tflog.Debug(ctx, "Retrieved group members from API state")

	var users []jiraGroupUsersModel
	for _, u := range members {
		m := &jiraGroupUsersModel{
			Self:         types.StringValue(u.Self),
			AccountID:    types.StringValue(u.AccountID),
			EmailAddress: types.StringValue(u.EmailAddress),
			AvatarUrls: &common.AvatarUrlsModel{
				One6X16:   types.StringValue(""),
				Two4X24:   types.StringValue(""),
				Three2X32: types.StringValue(""),
				Four8X48:  types.StringValue(""),
			},
			DisplayName: types.StringValue(u.DisplayName),
			Active:      types.BoolValue(u.Active),
			TimeZone:    types.StringValue(u.TimeZone),
			AccountType: types.StringValue(u.AccountType),
		}
		users = append(users, *m)
	}

	newState.ID = types.StringValue(group.Values[0].GroupID)
	newState.GroupID = types.StringValue(group.Values[0].GroupID)
	newState.Self = types.StringValue(fmt.Sprintf("https://%s/rest/api/3/group?groupId=%s", d.p.jira.Site.Host, group.Values[0].GroupID))
	newState.Users, _ = types.SetValueFrom(ctx, newState.Users.ElementType(ctx), users)

	tflog.Debug(ctx, "Storing group into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", newState),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
