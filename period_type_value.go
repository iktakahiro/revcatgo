package revcatgo

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

const (
	PeriodTypeTrial       = "TRIAL"
	PeriodTypeIntro       = "INTRO"
	PeriodTypeNormal      = "NORMAL"
	PeriodTypePromotional = "PROMOTIONAL"
)

var validPeriodTypeValues = []string{
	PeriodTypeTrial,
	PeriodTypeIntro,
	PeriodTypeNormal,
	PeriodTypePromotional,
}

type periodType struct {
	value null.String
}

func newPeriodType(s string) (*periodType, error) {
	if !contains(validPeriodTypeValues, s) {
		return &periodType{}, fmt.Errorf("periodType value should be one of the following: %v, got %v", strings.Join(validPeriodTypeValues, ","), s)
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
	_p, err := newPeriodType(strings.ToUpper(v.value.ValueOrZero()))
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of period_type: %w", err)
	}
	p.value = _p.value

	return nil
}
