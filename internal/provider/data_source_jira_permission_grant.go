package atlassian

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	jiraPermissionGrantDataSource struct {
		p atlassianProvider
	}

	jiraPermissionGrantDataSourceType struct{}

	jiraPermissionGrantDataSourceModel struct {
		ID                 types.String                    `tfsdk:"id"`
		PermissionSchemeID types.String                    `tfsdk:"permission_scheme_id"`
		Holder             *jiraPermissionGrantHolderModel `tfsdk:"holder"`
		Permission         types.String                    `tfsdk:"permission"`
	}
)

var (
	_ datasource.DataSource   = (*jiraPermissionGrantDataSource)(nil)
	_ provider.DataSourceType = (*jiraPermissionGrantDataSourceType)(nil)
)

func (d *jiraPermissionGrantDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *jiraPermissionGrantDataSourceType) NewDataSource(_ context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return &jiraPermissionGrantDataSource{
		p: provider,
	}, diags
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
