package int64modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Int64 = (*defaultValuePlanModifier)(nil)

type defaultValuePlanModifier struct {
	DefaultValue int64
}

func (m *defaultValuePlanModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m *defaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %q (%s)", m.DefaultValue, types.Int64Type)
}

func (m *defaultValuePlanModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// If the value is configured, skip validator
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// If the plan contains a value for the attribute, no need to proceed.
	// Do not override changes by a previous plan modifier.
	if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
		return
	}

	res.PlanValue = types.Int64Value(m.DefaultValue)
}

func DefaultValue(defaultValue int64) planmodifier.Int64 {
	return &defaultValuePlanModifier{
		DefaultValue: defaultValue,
	}
}
