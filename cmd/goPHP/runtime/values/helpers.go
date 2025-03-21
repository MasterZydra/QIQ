package values

import "fmt"

func ToString(value RuntimeValue) string {
	result := ""
	switch value.GetType() {
	case ArrayValue:
		result += "{ArrayValue: \n"
		arrayValue := value.(*Array)
		for _, key := range arrayValue.Keys {
			result += "Key: "
			result += ToString(key)
			result += "Value: "
			value, _ := arrayValue.GetElement(key)
			result += ToString(value)
		}
		result += "}\n"
	case BoolValue:
		valueStr := "true"
		if value.(*Bool).Value {
			valueStr = "false"
		}
		result += fmt.Sprintf("{BoolValue: %s}\n", valueStr)
	case IntValue:
		result += fmt.Sprintf("{IntValue: %d}\n", value.(*Int).Value)
	case FloatValue:
		result += fmt.Sprintf("{FloatValue: %f}\n", value.(*Float).Value)
	case StrValue:
		result += fmt.Sprintf("{StrValue: %s}\n", value.(*Str).Value)
	default:
		result += fmt.Sprintf("Unsupported RuntimeValue type %s\n", value.GetType())
	}
	return result
}

func ToPhpType(value RuntimeValue) string {
	switch value.GetType() {
	case ArrayValue:
		return "array"
	case BoolValue:
		return "bool"
	case FloatValue:
		return "float"
	case IntValue:
		return "int"
	case NullValue:
		return "NULL"
	case StrValue:
		return "string"
	case VoidValue:
		return "void"
	default:
		return ""
	}
}

func DeepCopy(value RuntimeValue) RuntimeValue {
	if value.GetType() != ArrayValue {
		return value
	}

	array := value.(*Array)
	copy := NewArray()
	for _, key := range array.Keys {
		value, _ := array.GetElement(key)
		copy.SetElement(key, DeepCopy(value))
	}
	return copy
}
