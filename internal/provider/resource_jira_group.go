package atlassian

import (
	"context"
	"fmt"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	common "github.com/openscientia/terraform-provider-atlassian/internal/provider/models"
)

type (
	jiraGroupResource struct {
		p atlassianProvider
	}

	jiraGroupResourceModel struct {
		ID      types.String `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		GroupID types.String `tfsdk:"group_id"`
		Self    types.String `tfsdk:"self"`
		Users   types.Set    `tfsdk:"users"`
	}

	jiraGroupUsersModel struct {
		Self         types.String            `tfsdk:"self"`
		AccountID    types.String            `tfsdk:"account_id"`
		EmailAddress types.String            `tfsdk:"email_address"`
		AvatarUrls   *common.AvatarUrlsModel `tfsdk:"avatar_urls"`
		DisplayName  types.String            `tfsdk:"display_name"`
		Active       types.Bool              `tfsdk:"active"`
		TimeZone     types.String            `tfsdk:"timezone"`
		AccountType  types.String            `tfsdk:"account_type"`
	}
)

var (
	_ resource.Resource                = (*jiraGroupResource)(nil)
	_ resource.ResourceWithImportState = (*jiraGroupResource)(nil)
)

func NewJiraGroupResource() resource.Resource {
	return &jiraGroupResource{}
}

func (*jiraGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_group"
}

func (*jiraGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version:             1,
		MarkdownDescription: "Jira Group Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group. Defaults to `group_id`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "(Forces new resource) The name of the group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group, which uniquely identifies the group across all Atlassian products.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"self": schema.StringAttribute{
				MarkdownDescription: "The URL for these group details.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
							Validators: []validator.String{
								stringvalidator.LengthAtMost(128),
							},
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *jiraGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jira.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jira.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.p.jira = client
}

func (*jiraGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *jiraGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating group resource")

	var plan jiraGroupResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded group plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	group, res, err := r.p.jira.Group.Create(ctx, plan.Name.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created group")

	bulkOptions := &models.GroupBulkOptionsScheme{
		GroupNames: []string{group.Name},
	}
	groupDetails, res, err := r.p.jira.Group.Bulk(ctx, bulkOptions, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve group details, got error: %s\n%s", err, resBody))
		return
	}

	plan.ID = types.StringValue(groupDetails.Values[0].GroupID)
	plan.GroupID = types.StringValue(groupDetails.Values[0].GroupID)
	plan.Self = types.StringValue(group.Self)
	plan.Users = types.SetNull(plan.Users.ElementType(ctx))

	tflog.Debug(ctx, "Storing group into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading group resource")

	var state jiraGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded group from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	bulkOptions := &models.GroupBulkOptionsScheme{
		GroupNames: []string{state.Name.ValueString()},
	}
	group, res, err := r.p.jira.Group.Bulk(ctx, bulkOptions, 0, 1)
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get group, got error: %s\n%s", err, resBody))
		return
	}

	isLast := false
	startAt := 0
	maxResults := 100
	members := []*models.GroupUserDetailScheme{}
	for !isLast {
		groupMembers, res, err := r.p.jira.Group.Members(ctx, state.Name.ValueString(), true, startAt, maxResults)
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

	tflog.Debug(ctx, "Retrieved group from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v, Members Count:%+v", group.Values[0], len(members)),
	})

	state.ID = types.StringValue(group.Values[0].GroupID)
	state.GroupID = types.StringValue(group.Values[0].GroupID)
	state.Self = types.StringValue(fmt.Sprintf("https://%s/rest/api/3/group?groupId=%s", r.p.jira.Site.Host, group.Values[0].GroupID))

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
	state.Users, _ = types.SetValueFrom(ctx, state.Users.ElementType(ctx), users)

	tflog.Debug(ctx, "Storing group into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The RequiresReplace plan modifier will trigger Terraform to destroy and recreate the resource
	// if any of the required attributes changes, i.e. name.
	tflog.Debug(ctx, "If the value of any required attribute changes, Terraform will destroy and recreate the resource")
}

func (r *jiraGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting group resource")

	var state jiraGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.p.jira.Group.Delete(ctx, state.Name.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete group, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted group from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
