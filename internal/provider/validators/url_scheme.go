package validators

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ validator.String = (*urlWithSchemeValidator)(nil)

type urlWithSchemeValidator struct {
	values []string
}

func (v urlWithSchemeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v urlWithSchemeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Must be a URL and its scheme is one of: %q", v.values)
}

func (v urlWithSchemeValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	tflog.Debug(ctx, "Validating attribute value is a URL with acceptable scheme", map[string]interface{}{
		"attribute":         req.Path.String(),
		"acceptableSchemes": strings.Join(v.values, ","),
	})

	var val types.String
	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &val)
	if diags.HasError() {
		res.Diagnostics.Append(diags...)
	}

	if val.IsNull() || val.IsUnknown() {
		return
	}

	u, err := url.Parse(val.ValueString())
	if err != nil {
		res.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			fmt.Sprintf("Parsing URL %q failed: %v", val.ValueString(), err),
		)
		return
	}

	if u.Host == "" {
		res.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			fmt.Sprintf("URL %q contains no host", u.String()),
		)
		return
	}

	for _, s := range v.values {
		if u.Scheme == s {
			return
		}
	}

	res.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid URL scheme",
		fmt.Sprintf("URL %q expected to use scheme from %q, got: %q", u.String(), v.values, u.Scheme),
	)
}

func UrlWithScheme(values ...string) validator.String {
	return urlWithSchemeValidator{
		values: values,
	}
}
