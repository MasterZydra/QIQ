package values

type ValueType string

const (
	VoidValue  ValueType = "Void"
	NullValue  ValueType = "Null"
	ArrayValue ValueType = "Array"
	BoolValue  ValueType = "Bool"
	IntValue   ValueType = "Int"
	FloatValue ValueType = "Float"
	StrValue   ValueType = "Str"
)
