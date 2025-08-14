package variableHandling

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime/values"
)

func ToValueType(valueType values.ValueType, runtimeValue values.RuntimeValue, leadingNumeric bool) (values.RuntimeValue, phpError.Error) {
	switch valueType {
	case values.BoolValue:
		boolean, err := BoolVal(runtimeValue)
		return values.NewBool(boolean), err
	case values.FloatValue:
		floating, err := FloatVal(runtimeValue, leadingNumeric)
		return values.NewFloat(floating), err
	case values.IntValue:
		integer, err := IntVal(runtimeValue, leadingNumeric)
		return values.NewInt(integer), err
	case values.StrValue:
		str, err := StrVal(runtimeValue)
		return values.NewStr(str), err
	default:
		return values.NewVoid(), phpError.NewError("runtimeValueToValueType: Unsupported runtime value: %s", valueType)
	}
}
