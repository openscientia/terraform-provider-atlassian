package common

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AvatarUrlsModel struct {
	One6X16   types.String `tfsdk:"p16x16"`
	Two4X24   types.String `tfsdk:"p24x24"`
	Three2X32 types.String `tfsdk:"p32x32"`
	Four8X48  types.String `tfsdk:"p48x48"`
}
