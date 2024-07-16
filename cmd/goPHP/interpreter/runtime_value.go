package interpreter

type ValueType string

const (
	VoidValue     ValueType = "Void"
	NullValue     ValueType = "Null"
	ArrayValue    ValueType = "Array"
	BooleanValue  ValueType = "Boolean"
	IntegerValue  ValueType = "Integer"
	FloatingValue ValueType = "Floating"
	StringValue   ValueType = "String"
)

// MARK: RuntimeValue

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

// MARK: VoidValue

func NewVoidRuntimeValue() *RuntimeValue {
	return &RuntimeValue{valueType: VoidValue}
}

// MARK: NullValue

func NewNullRuntimeValue() *RuntimeValue {
	return &RuntimeValue{valueType: NullValue}
}

// MARK: ArrayRuntimeValue

type ArrayRuntimeValue struct {
	*RuntimeValue
	Keys     []IRuntimeValue
	Elements map[IRuntimeValue]IRuntimeValue
}

func NewArrayRuntimeValue() *ArrayRuntimeValue {
	return &ArrayRuntimeValue{
		RuntimeValue: NewRuntimeValue(ArrayValue),
		Keys:         []IRuntimeValue{},
		Elements:     map[IRuntimeValue]IRuntimeValue{},
	}
}

func NewArrayRuntimeValueFromMap(elements map[IRuntimeValue]IRuntimeValue) *ArrayRuntimeValue {
	keys := []IRuntimeValue{}
	for key := range elements {
		keys = append(keys, key)
	}
	return &ArrayRuntimeValue{
		RuntimeValue: NewRuntimeValue(ArrayValue),
		Keys:         keys,
		Elements:     elements,
	}
}

func (runtimeValue *ArrayRuntimeValue) SetElement(key IRuntimeValue, value IRuntimeValue) {
	existingKey, exists := runtimeValue.findKey(key)
	if !exists {
		runtimeValue.Keys = append(runtimeValue.Keys, key)
		runtimeValue.Elements[key] = value
	} else {
		runtimeValue.Elements[existingKey] = value
	}
}

func (runtimeValue *ArrayRuntimeValue) findKey(key IRuntimeValue) (IRuntimeValue, bool) {
	for k := range runtimeValue.Elements {
		if k.GetType() != key.GetType() {
			continue
		}
		boolean, err := compare(key, "===", k)
		if err != nil {
			return NewVoidRuntimeValue(), false
		}
		if boolean.Value {
			return k, true
		}
	}
	return NewVoidRuntimeValue(), false
}

func (runtimeValue *ArrayRuntimeValue) GetElement(key IRuntimeValue) (IRuntimeValue, bool) {
	key, found := runtimeValue.findKey(key)
	if !found {
		return NewVoidRuntimeValue(), false
	}
	return runtimeValue.Elements[key], true
}

// MARK: BooleanRuntimeValue

type BooleanRuntimeValue struct {
	*RuntimeValue
	Value bool
}

func NewBooleanRuntimeValue(value bool) *BooleanRuntimeValue {
	return &BooleanRuntimeValue{RuntimeValue: NewRuntimeValue(BooleanValue), Value: value}
}

// MARK: IntegerRuntimeValue

type IntegerRuntimeValue struct {
	*RuntimeValue
	Value int64
}

func NewIntegerRuntimeValue(value int64) *IntegerRuntimeValue {
	return &IntegerRuntimeValue{RuntimeValue: NewRuntimeValue(IntegerValue), Value: value}
}

// MARK: FloatingRuntimeValue

type FloatingRuntimeValue struct {
	*RuntimeValue
	Value float64
}

func NewFloatingRuntimeValue(value float64) *FloatingRuntimeValue {
	return &FloatingRuntimeValue{RuntimeValue: NewRuntimeValue(FloatingValue), Value: value}
}

// MARK: StringRuntimeValue

type StringRuntimeValue struct {
	*RuntimeValue
	Value string
}

func NewStringRuntimeValue(value string) *StringRuntimeValue {
	return &StringRuntimeValue{RuntimeValue: NewRuntimeValue(StringValue), Value: value}
}
