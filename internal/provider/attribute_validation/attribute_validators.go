package attribute_validation

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type intValuesValidator struct {
	Values []int
}

func (v intValuesValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("only valid int values are %v", v.Values)
}

func (v intValuesValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("only valid int values are %v", v.Values)
}

func (v intValuesValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {

	var num types.Int64
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &num)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if num.Unknown || num.Null {
		return
	}

	flag := false
	for _, s := range v.Values {
		if num.Value == int64(s) {
			flag = true
		}
	}

	if !flag {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Value",
			fmt.Sprintf("only valid int values are %v", v.Values),
		)

		return
	}
}

func IntValues(values []int) intValuesValidator {
	return intValuesValidator{
		Values: values,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type stringLengthBetweenValidator struct {
	Min int
	Max int
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLengthBetweenValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", v.Min, v.Max)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string length must be between `%d` and `%d`", v.Min, v.Max)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringLengthBetweenValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic validators, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown || str.Null {
		return
	}

	strLen := len(str.Value)

	if strLen < v.Min || strLen > v.Max {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Length",
			fmt.Sprintf("String length must be between %d and %d, got: %d.", v.Min, v.Max, strLen),
		)

		return
	}
}

func StringLengthBetween(minLength int, maxLength int) stringLengthBetweenValidator {
	return stringLengthBetweenValidator{
		Max: maxLength,
		Min: minLength,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type stringValuesValidator struct {
	Values []string
}

func (v stringValuesValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("only valid string values are `%v`", v.Values)
}

func (v stringValuesValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("only valid string values are `%v`", v.Values)
}

func (v stringValuesValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {

	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown || str.Null {
		return
	}

	flag := false
	for _, s := range v.Values {
		if str.Value == s {
			flag = true
		}
	}

	if !flag {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Value",
			fmt.Sprintf("only valid string values are `%v`", v.Values),
		)

		return
	}
}

func StringValues(values []string) stringValuesValidator {
	return stringValuesValidator{
		Values: values,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type urlWithSchemeAttributeValidator struct {
	acceptableSchemes []string
}

func UrlWithScheme(acceptableSchemes ...string) tfsdk.AttributeValidator {
	return &urlWithSchemeAttributeValidator{acceptableSchemes}
}

var _ tfsdk.AttributeValidator = (*urlWithSchemeAttributeValidator)(nil)

func (v *urlWithSchemeAttributeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v *urlWithSchemeAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Must be a URL and its scheme is one of: %q", v.acceptableSchemes)
}

func (v *urlWithSchemeAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	tflog.Debug(ctx, "Validating attribute value is a URL with acceptable scheme", map[string]interface{}{
		"attribute":         req.AttributePath.String(),
		"acceptableSchemes": strings.Join(v.acceptableSchemes, ","),
	})

	var val types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &val)
	if diags.HasError() {
		res.Diagnostics.Append(diags...)
	}

	if val.IsNull() || val.IsUnknown() {
		return
	}

	u, err := url.Parse(val.Value)
	if err != nil {
		res.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid URL",
			fmt.Sprintf("Parsing URL %q failed: %v", val.Value, err),
		)
		return
	}

	if u.Host == "" {
		res.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid URL",
			fmt.Sprintf("URL %q contains no host", u.String()),
		)
		return
	}

	for _, s := range v.acceptableSchemes {
		if u.Scheme == s {
			return
		}
	}

	res.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid URL scheme",
		fmt.Sprintf("URL %q expected to use scheme from %q, got: %q", u.String(), v.acceptableSchemes, u.Scheme),
	)
}
