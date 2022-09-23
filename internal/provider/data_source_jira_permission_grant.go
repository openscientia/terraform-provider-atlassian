package atlassian

import (
	"context"
	"fmt"
	"strconv"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraPermissionGrantDataSource struct {
		p atlassianProvider
	}

	jiraPermissionGrantDataSourceModel struct {
		ID                 types.String                    `tfsdk:"id"`
		PermissionSchemeID types.String                    `tfsdk:"permission_scheme_id"`
		Holder             *jiraPermissionGrantHolderModel `tfsdk:"holder"`
		Permission         types.String                    `tfsdk:"permission"`
	}
)

var (
	_ datasource.DataSource = (*jiraPermissionGrantDataSource)(nil)
)

func NewJiraPermissionGrantDataSource() datasource.DataSource {
	return &jiraPermissionGrantDataSource{}
}

func (*jiraPermissionGrantDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jira_permission_grant"
}

func (*jiraPermissionGrantDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version:             1,
		MarkdownDescription: "Jira Permission Grant Data Source",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "The ID of the permission grant.",
				Required:            true,
				Type:                types.StringType,
			},
			"permission_scheme_id": {
				MarkdownDescription: "The ID of the permission scheme in which to create a new permission grant.",
				Required:            true,
				Type:                types.StringType,
			},
			"holder": {
				MarkdownDescription: "The user, group, field or role being granted the permission.",
				Computed:            true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"type": {
							MarkdownDescription: "The type of permission holder.",
							Computed:            true,
							Type:                types.StringType,
						},
						"parameter": {
							MarkdownDescription: "The identifier associated with the `type` value that defines the holder of the permission.",
							Computed:            true,
							Type:                types.StringType,
						},
					},
				),
			},
			"permission": {
				MarkdownDescription: "The permission to grant.",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (d *jiraPermissionGrantDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jiraPermissionGrantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading permission grant data source")

	var newState jiraPermissionGrantDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Loaded permission grant config", map[string]interface{}{
		"readConfig": fmt.Sprintf("%+v", newState),
	})

	grantId, _ := strconv.Atoi(newState.ID.Value)
	schemeId, _ := strconv.Atoi(newState.PermissionSchemeID.Value)
	permissionGrant, res, err := d.p.jira.Permission.Scheme.Grant.Get(ctx, schemeId, grantId, []string{"all"})
	if err != nil {
		var resBody string
		if res != nil {
			resBody = res.Bytes.String()
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get permission grant, got error: %s\n%s", err, resBody))
		return
	}
	tflog.Debug(ctx, "Retrieved permission grant from API state", map[string]interface{}{
		"readApiState": fmt.Sprintf("%+v, Holder:%+v", permissionGrant, permissionGrant.Holder),
	})

	newState.Holder = &jiraPermissionGrantHolderModel{
		Type:      types.String{Value: permissionGrant.Holder.Type},
		Parameter: types.String{Value: permissionGrant.Holder.Parameter},
	}
	newState.Permission = types.String{Value: permissionGrant.Permission}

	tflog.Debug(ctx, "Storing permission grant into the state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
