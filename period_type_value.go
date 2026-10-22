package revcatgo

import (
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

const (
	// PeriodTypeTrial identifies a trial period.
	PeriodTypeTrial = "TRIAL"
	// PeriodTypeIntro identifies an introductory pricing period.
	PeriodTypeIntro = "INTRO"
	// PeriodTypeNormal identifies a standard billing period.
	PeriodTypeNormal = "NORMAL"
	// PeriodTypePromotional identifies a promotional period.
	PeriodTypePromotional = "PROMOTIONAL"
	// PeriodTypePrepaid identifies a prepaid subscription period.
	PeriodTypePrepaid = "PREPAID"
)

var validPeriodTypeValues = []string{
	PeriodTypeTrial,
	PeriodTypeIntro,
	PeriodTypeNormal,
	PeriodTypePromotional,
	PeriodTypePrepaid,
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

// IsKnown reports whether RevenueCat documents the period type as one of the PeriodType constants.
// Period types added by RevenueCat later (and null values) return false but remain decodable.
func (p periodType) IsKnown() bool {
	return contains(validPeriodTypeValues, p.value.ValueOrZero())
}

func (p periodType) MarshalJSON() ([]byte, error) {
	return p.value.MarshalJSON()
}

// UnmarshalJSON deserializes a period type from JSON.
// JSON null and empty strings decode to the zero value, and undocumented values are
// preserved so that new RevenueCat values do not break decoding; see IsKnown.
func (p *periodType) UnmarshalJSON(b []byte) error {
	v, err := unmarshalEnumValue(b, true)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of period_type: %w", err)
	}
	p.value = v

	return nil
}
