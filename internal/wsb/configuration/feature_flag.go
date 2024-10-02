package configuration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type FeatureFlag int

const (
	NotSet FeatureFlag = iota
	Enable
	Disable
	Default
)

func (f FeatureFlag) String() string {
	return [...]string{"", "Enable", "Disable", "Default"}[f]
}

func (f FeatureFlag) ToBool(defaultValue bool) basetypes.BoolValue {
	switch f {
	case Enable:
		return types.BoolValue(true)
	case Disable:
		return types.BoolValue(false)
	case Default:
		return types.BoolValue(defaultValue)
	default:
		return types.BoolNull()
	}
}

func ParseFeatureFlag(s string) (FeatureFlag, error) {
	switch strings.ToLower(s) {
	case "enable":
		return Enable, nil
	case "disable":
		return Disable, nil
	case "default":
		return Default, nil
	case "":
		return NotSet, nil
	default:
		return Default, fmt.Errorf("unknown feature flag: %s", s)
	}
}

func ToFeatureFlag(value bool) FeatureFlag {
	switch value {
	case true:
		return Enable
	default:
		return Disable
	}
}
