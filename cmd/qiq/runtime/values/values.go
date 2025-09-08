package values

import (
	"strconv"
)

// MARK: RuntimeValue

type RuntimeValue interface {
	GetType() ValueType
}

// MARK: abstractValue

type abstractValue struct {
	valueType ValueType
}

func newAbstractValue(valueType ValueType) *abstractValue {
	return &abstractValue{valueType: valueType}
}

func (value *abstractValue) GetType() ValueType {
	return value.valueType
}

// MARK: Void

type Void struct {
	*abstractValue
}

var void = &Void{abstractValue: newAbstractValue(VoidValue)}

func NewVoid() *Void {
	return void
}

func NewVoidSlot() *Slot {
	return NewSlot(NewVoid())
}

// MARK: Null

type Null struct {
	*abstractValue
}

var null = &Null{abstractValue: newAbstractValue(NullValue)}

func NewNull() *Null {
	return null
}

func NewNullSlot() *Slot {
	return NewSlot(NewNull())
}

// MARK: Bool

type Bool struct {
	*abstractValue
	Value bool
}

var trueRuntimeValue = &Bool{abstractValue: newAbstractValue(BoolValue), Value: true}
var falseRuntimeValue = &Bool{abstractValue: newAbstractValue(BoolValue), Value: false}

func NewBool(value bool) *Bool {
	if value {
		return trueRuntimeValue
	}
	return falseRuntimeValue
}

func NewBoolSlot(value bool) *Slot {
	return NewSlot(NewBool(value))
}

// MARK: Int

type Int struct {
	*abstractValue
	Value int64
}

func NewInt(value int64) *Int {
	return &Int{abstractValue: newAbstractValue(IntValue), Value: value}
}

func NewIntSlot(value int64) *Slot {
	return NewSlot(NewInt(value))
}

// MARK: Float

type Float struct {
	*abstractValue
	Value float64
}

func NewFloat(value float64) *Float {
	return &Float{abstractValue: newAbstractValue(FloatValue), Value: value}
}

func NewFloatSlot(value float64) *Slot {
	return NewSlot(NewFloat(value))
}

func (value *Float) ToPhpString() string {
	return strconv.FormatFloat(value.Value, 'f', -1, 64)
}

// MARK: Str

type Str struct {
	*abstractValue
	Value string
}

func NewStr(value string) *Str {
	return &Str{abstractValue: newAbstractValue(StrValue), Value: value}
}

func NewStrSlot(value string) *Slot {
	return NewSlot(NewStr(value))
}
