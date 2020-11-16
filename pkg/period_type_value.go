package revcatgo

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

const (
	periodTypeTrial       = "TRIAL"
	periodTypeIntro       = "INTRO"
	periodTypeNormal      = "NORMAL"
	periodTypePromotional = "PROMOTIONAL"
)

var validperiodTypeValues = []string{
	periodTypeTrial,
	periodTypeIntro,
	periodTypeNormal,
	periodTypePromotional,
}

type periodType struct {
	value null.String
}

func newPeriodType(s string) (*periodType, error) {
	if !contains(validperiodTypeValues, s) {
		return &periodType{}, errors.New("periodType value should be one of the following:" + strings.Join(validperiodTypeValues, ","))
	}
	return &periodType{value: null.StringFrom(s)}, nil
}

func (p periodType) String() string {
	return p.value.ValueOrZero()
}

func (p periodType) MarshalJSON() ([]byte, error) {
	return p.value.MarshalJSON()
}

// UnmarshalJSON deserializes a store from JSON
func (p *periodType) UnmarshalJSON(b []byte) error {
	v := &periodType{}
	err := v.value.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of period_type: %w", err)
	}
	if !v.value.Valid {
		return errors.New("period_type is a required field")
	}
	_p, err := newPeriodType(v.value.ValueOrZero())
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of period_type: %w", err)
	}
	p.value = _p.value

	return nil
}
