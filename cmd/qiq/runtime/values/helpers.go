package values

import "fmt"

func ToString(value RuntimeValue) string {
	switch value.GetType() {
	case ArrayValue:
		result := "{ArrayValue: \n"
		arrayValue := value.(*Array)
		for _, key := range arrayValue.Keys {
			result += "Key: "
			result += ToString(key)
			result += "Value: "
			value, _ := arrayValue.GetElement(key)
			result += ToString(value.Value)
		}
		result += "}\n"
		return result
	case BoolValue:
		valueStr := "true"
		if value.(*Bool).Value {
			valueStr = "false"
		}
		return fmt.Sprintf("{BoolValue: %s}\n", valueStr)
	case IntValue:
		return fmt.Sprintf("{IntValue: %d}\n", value.(*Int).Value)
	case FloatValue:
		return fmt.Sprintf("{FloatValue: %f}\n", value.(*Float).Value)
	case StrValue:
		return fmt.Sprintf("{StrValue: %s}\n", value.(*Str).Value)
	case NullValue:
		return "{NullValue}"
	case VoidValue:
		return "{VoidValue}"
	case ObjectValue:
		return fmt.Sprintf("{Object: %s}\n", value.(*Object).Class.Name)
		// TODO Add properties
	default:
		return fmt.Sprintf("Unsupported RuntimeValue type %s\n", value.GetType())
	}
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
	case ObjectValue:
		return "object"
	case VoidValue:
		return "void"
	default:
		return ""
	}
}

func DeepCopy(slot *Slot) *Slot {
	value := slot.Value
	if value.GetType() != ArrayValue {
		return NewSlot(slot.Value)
	}

	array := value.(*Array)
	copy := NewArray()
	for _, key := range array.Keys {
		value, _ := array.GetElement(key)
		copy.SetElement(key, DeepCopy(value).Value)
	}
	return NewSlot(copy)
}
