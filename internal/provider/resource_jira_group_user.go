package atlassian

import (
	"context"
	"fmt"
	"strings"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	common "github.com/openscientia/terraform-provider-atlassian/internal/provider/models"
)

type (
	jiraGroupUserResource struct {
		p atlassianProvider
	}

	jiraGroupUserResourceModel struct {
		ID           types.String            `tfsdk:"id"`
		GroupName    types.String            `tfsdk:"group_name"`
		AccountID    types.String            `tfsdk:"account_id"`
		Self         types.String            `tfsdk:"self"`
		EmailAddress types.String            `tfsdk:"email_address"`
		AvatarUrls   *common.AvatarUrlsModel `tfsdk:"avatar_urls"`
		DisplayName  types.String            `tfsdk:"display_name"`
		Active       types.Bool              `tfsdk:"active"`
		TimeZone     types.String            `tfsdk:"timezone"`
		AccountType  types.String            `tfsdk:"account_type"`
	}
)

var (
	_ resource.Resource                = (*jiraGroupUserResource)(nil)
	_ resource.ResourceWithImportState = (*jiraGroupUserResource)(nil)
)

func NewJiraGroupUserResource() resource.Resource {
	return &jiraGroupUserResource{}
}

func (*jiraGroupUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_group_user"
}

func (*jiraGroupUserResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Group User Resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the group user. It is computed using `group_name` and `account_id` separated by a hyphen (`-`).",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"group_name": {
				MarkdownDescription: "(Forces new resource) The name of the group.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"account_id": {
				MarkdownDescription: "(Forces new resource) The account ID of the user, which uniquely identifies the user across all Atlassian products.",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.RequiresReplace(),
				},
			},
			"self": {
				MarkdownDescription: "The URL of the user.",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"email_address": {
				MarkdownDescription: "The email address of the user. Depending on the user’s privacy settings, this may be returned as null.",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"avatar_urls": {
				MarkdownDescription: "The avatars of the user.",
				Computed:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"p16x16": {
						MarkdownDescription: "The URL of the item's 16x16 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"p24x24": {
						MarkdownDescription: "The URL of the item's 24x24 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"p32x32": {
						MarkdownDescription: "The URL of the item's 32x32 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
					"p48x48": {
						MarkdownDescription: "The URL of the item's 48x48 pixel avatar.",
						Computed:            true,
						Type:                types.StringType,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
					},
				}),
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"display_name": {
				MarkdownDescription: "The display name of the user. Depending on the user’s privacy settings, this may return an alternative value.",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"active": {
				MarkdownDescription: "Whether the user is active.",
				Computed:            true,
				Type:                types.BoolType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"timezone": {
				MarkdownDescription: "The time zone specified in the user's profile. Depending on the user’s privacy settings, this may be returned as null.",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"account_type": {
				MarkdownDescription: "The type of account represented by this user. This will be one of `atlassian` (normal users), `app` (application user) or `customer` (Jira Service Desk customer user).",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
		}}, nil
}

func (r *jiraGroupUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (*jiraGroupUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: group_name, account_id. Got: %q", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_id"), idParts[1])...)
}

func (r *jiraGroupUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating group user resource")

	var plan jiraGroupUserResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded group user plan", map[string]interface{}{
		"createPlan": fmt.Sprintf("%+v", plan),
	})

	_, res, err := r.p.jira.Group.Add(ctx, plan.GroupName.ValueString(), plan.AccountID.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create group user, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Created group user")

	isLast := false
	startAt := 0
	maxResults := 100
	users := []*models.GroupUserDetailScheme{}
	for !isLast {
		groupUsers, res, err := r.p.jira.Group.Members(ctx, plan.GroupName.ValueString(), true, startAt, maxResults)
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get group users, got error: %s\n%s", err, resBody))
			return
		}
		startAt += maxResults
		isLast = groupUsers.IsLast
		users = append(users, groupUsers.Values...)
	}
	tflog.Debug(ctx, "Retrieved group users from API state")

	for _, u := range users {
		if u.AccountID == plan.AccountID.ValueString() {
			plan.Self = types.StringValue(u.Self)
			plan.EmailAddress = types.StringValue(u.EmailAddress)
			plan.AvatarUrls = &common.AvatarUrlsModel{
				One6X16:   types.StringValue(""),
				Two4X24:   types.StringValue(""),
				Three2X32: types.StringValue(""),
				Four8X48:  types.StringValue(""),
			}
			plan.DisplayName = types.StringValue(u.DisplayName)
			plan.Active = types.BoolValue(u.Active)
			plan.TimeZone = types.StringValue(u.TimeZone)
			plan.AccountType = types.StringValue(u.AccountType)
			continue
		}
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s-%s", plan.GroupName.ValueString(), plan.AccountID.ValueString()))

	tflog.Debug(ctx, "Storing group user into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", plan),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jiraGroupUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading group user resource")

	var state jiraGroupUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded group user from state", map[string]interface{}{
		"readState": fmt.Sprintf("%+v", state),
	})

	isLast := false
	startAt := 0
	maxResults := 100
	users := []*models.GroupUserDetailScheme{}
	for !isLast {
		groupUsers, res, err := r.p.jira.Group.Members(ctx, state.GroupName.ValueString(), true, startAt, maxResults)
		if err != nil {
			var resBody string
			if res != nil {
				resBody = res.Bytes.String()
			}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get group users, got error: %s\n%s", err, resBody))
			return
		}
		startAt += maxResults
		isLast = groupUsers.IsLast
		users = append(users, groupUsers.Values...)
	}
	tflog.Debug(ctx, "Retrieved group users from API state")

	for _, u := range users {
		if u.AccountID == state.AccountID.ValueString() {
			state.Self = types.StringValue(u.Self)
			state.EmailAddress = types.StringValue(u.EmailAddress)
			state.AvatarUrls = &common.AvatarUrlsModel{
				One6X16:   types.StringValue(""),
				Two4X24:   types.StringValue(""),
				Three2X32: types.StringValue(""),
				Four8X48:  types.StringValue(""),
			}
			state.DisplayName = types.StringValue(u.DisplayName)
			state.Active = types.BoolValue(u.Active)
			state.TimeZone = types.StringValue(u.TimeZone)
			state.AccountType = types.StringValue(u.AccountType)
			continue
		}
	}
	state.ID = types.StringValue(fmt.Sprintf("%s-%s", state.GroupName.ValueString(), state.AccountID.ValueString()))

	tflog.Debug(ctx, "Storing group user into the state", map[string]interface{}{
		"readNewState": fmt.Sprintf("%+v", state),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jiraGroupUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The RequiresReplace plan modifier will trigger Terraform to destroy and recreate the resource
	// if any of the required attributes changes, i.e. group_name and/or account_id.
	tflog.Debug(ctx, "If the value of any required attribute changes, Terraform will destroy and recreate the resource")
}

func (r *jiraGroupUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting group user resource")

	var state jiraGroupUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.p.jira.Group.Remove(ctx, state.GroupName.ValueString(), state.AccountID.ValueString())
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete group user, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Deleted group user from API state")

	// If a Resource type Delete method is completed without error, the framework will automatically remove the resource.
}
