package interpreter

type ValueType string

const (
	VoidValue     ValueType = "Void"
	NullValue     ValueType = "Null"
	IntegerValue  ValueType = "Integer"
	FloatingValue ValueType = "Floating"
	StringValue   ValueType = "String"
)

// RuntimeValue

type IRuntimeValue interface {
	GetType() ValueType
}

type RuntimeValue struct {
	valueType ValueType
}

func NewRuntimeValue(valueType ValueType) *RuntimeValue {
	return &RuntimeValue{valueType: valueType}
}

func (runtimeValue *RuntimeValue) GetType() ValueType {
	return runtimeValue.valueType
}

// VoidValue

func NewVoidRuntimeValue() *RuntimeValue {
	return &RuntimeValue{valueType: VoidValue}
}

// NullValue

func NewNullRuntimeValue() *RuntimeValue {
	return &RuntimeValue{valueType: NullValue}
}

// IntegerRuntimeValue

type IIntegerRuntimeValue interface {
	IRuntimeValue
	GetValue() int64
}

type IntegerRuntimeValue struct {
	runtimeValue *RuntimeValue
	value        int64
}

func NewIntegerRuntimeValue(value int64) *IntegerRuntimeValue {
	return &IntegerRuntimeValue{runtimeValue: NewRuntimeValue(IntegerValue), value: value}
}

func (runtimeValue *IntegerRuntimeValue) GetType() ValueType {
	return runtimeValue.runtimeValue.valueType
}

func (runtimeValue *IntegerRuntimeValue) GetValue() int64 {
	return runtimeValue.value
}

func runtimeValToIntRuntimeVal(runtimeValue IRuntimeValue) IIntegerRuntimeValue {
	var i interface{} = runtimeValue
	return i.(IIntegerRuntimeValue)
}

// FloatingRuntimeValue

type IFloatingRuntimeValue interface {
	IRuntimeValue
	GetValue() float64
}

type FloatingRuntimeValue struct {
	runtimeValue *RuntimeValue
	value        float64
}

func NewFloatingRuntimeValue(value float64) *FloatingRuntimeValue {
	return &FloatingRuntimeValue{runtimeValue: NewRuntimeValue(FloatingValue), value: value}
}

func (runtimeValue *FloatingRuntimeValue) GetType() ValueType {
	return runtimeValue.runtimeValue.valueType
}

func (runtimeValue *FloatingRuntimeValue) GetValue() float64 {
	return runtimeValue.value
}

func runtimeValToFloatRuntimeVal(runtimeValue IRuntimeValue) IFloatingRuntimeValue {
	var i interface{} = runtimeValue
	return i.(IFloatingRuntimeValue)
}

// StringRuntimeValue

type IStringRuntimeValue interface {
	IRuntimeValue
	GetValue() string
}

type StringRuntimeValue struct {
	runtimeValue *RuntimeValue
	value        string
}

func NewStringRuntimeValue(value string) *StringRuntimeValue {
	return &StringRuntimeValue{runtimeValue: NewRuntimeValue(StringValue), value: value}
}

func (runtimeValue *StringRuntimeValue) GetType() ValueType {
	return runtimeValue.runtimeValue.valueType
}

func (runtimeValue *StringRuntimeValue) GetValue() string {
	return runtimeValue.value
}

func runtimeValToStrRuntimeVal(runtimeValue IRuntimeValue) IStringRuntimeValue {
	var i interface{} = runtimeValue
	return i.(IStringRuntimeValue)
}
