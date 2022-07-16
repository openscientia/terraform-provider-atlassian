package attribute_validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
