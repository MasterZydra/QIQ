package values

type Slot struct {
	Value RuntimeValue
}

func NewSlot(value RuntimeValue) *Slot { return &Slot{Value: value} }

func (slot Slot) GetType() ValueType { return slot.Value.GetType() }
