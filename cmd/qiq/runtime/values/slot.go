package values

type Slot struct {
	Value RuntimeValue
}

func NewSlot(value RuntimeValue) *Slot {
	return &Slot{Value: value}
}
