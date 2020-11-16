package revcatgo

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type price struct {
	value null.Float
}

func newPrice(v null.Float) (*price, error) {
	return &price{value: v}, nil
}

// String returns the price value.
func (p price) Float64() float64 {
	return p.value.ValueOrZero()
}

func (p price) IsFreeTrial() bool {
	return p.value.Float64 == 0
}

func (p price) MarshalJSON() ([]byte, error) {
	return p.value.MarshalJSON()
}

// UnmarshalJSON deserializes a store from JSON
func (p *price) UnmarshalJSON(b []byte) error {
	v := &price{}
	err := v.value.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of price: %w", err)
	}
	_p, err := newPrice(v.value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of price: %w", err)
	}
	p.value = _p.value

	return nil
}
