package stringmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = (*defaultValuePlanModifier)(nil)

type defaultValuePlanModifier struct {
	DefaultValue string
}

func (m *defaultValuePlanModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m *defaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %q (%s)", m.DefaultValue, types.StringType)
}

func (m *defaultValuePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	// If the value is configured, skip validator
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// If the plan contains a value for the attribute, no need to proceed.
	// Do not override changes by a previous plan modifier.
	if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
		return
	}

	res.PlanValue = types.StringValue(m.DefaultValue)
}

func DefaultValue(defaultValue string) planmodifier.String {
	return &defaultValuePlanModifier{
		DefaultValue: defaultValue,
	}
}
