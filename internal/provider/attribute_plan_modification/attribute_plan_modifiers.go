package attribute_plan_modification

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type stringDefaultModifier struct {
	Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// Modify runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m stringDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic plan modifiers, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributePlan, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If value is configured, skip validator
	if !str.IsNull() && !str.IsUnknown() {
		return
	}

	resp.AttributePlan = types.String{Value: m.Default}
}

func StringDefault(defaultValue string) stringDefaultModifier {
	return stringDefaultModifier{
		Default: defaultValue,
	}
}

/* -------------------------------------------------------------------------------------------------------------------------------------------- */

type defaultValueAttributePlanModifier struct {
	DefaultValue attr.Value
}

func DefaultValue(v attr.Value) tfsdk.AttributePlanModifier {
	return &defaultValueAttributePlanModifier{v}
}

var _ tfsdk.AttributePlanModifier = (*defaultValueAttributePlanModifier)(nil)

func (m *defaultValueAttributePlanModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m *defaultValueAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %q (%s)", m.DefaultValue, m.DefaultValue.Type(ctx))
}

func (m *defaultValueAttributePlanModifier) Modify(_ context.Context, req tfsdk.ModifyAttributePlanRequest, res *tfsdk.ModifyAttributePlanResponse) {
	// If the value is configured, skip validator
	if !req.AttributeConfig.IsNull() && !req.AttributeConfig.IsUnknown() {
		return
	}

	// If the plan contains a value for the attribute, no need to proceed.
	// Do not override changes by a previous plan modifier.
	if !req.AttributePlan.IsNull() && !req.AttributePlan.IsUnknown() {
		return
	}

	res.AttributePlan = m.DefaultValue
}
