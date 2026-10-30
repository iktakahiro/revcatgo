package revcatgo

import (
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type environment struct {
	value null.String
}

const (
	// EnvironmentSandbox marks events that occurred in the sandbox environment.
	EnvironmentSandbox = "SANDBOX"
	// EnvironmentProduction marks events that occurred in the production environment.
	EnvironmentProduction = "PRODUCTION"
)

var validEnvironmentValues = []string{
	EnvironmentSandbox,
	EnvironmentProduction,
}

func newEnvironment(s string) (*environment, error) {
	if !contains(validEnvironmentValues, s) {
		return &environment{}, fmt.Errorf("environment value should be one of the following: %v, got %v", strings.Join(validEnvironmentValues, ","), s)
	}
	return &environment{value: null.StringFrom(s)}, nil
}

func (e environment) String() string {
	return e.value.ValueOrZero()
}

func (e *environment) IsProduction() bool {
	return e.value.ValueOrZero() == EnvironmentProduction
}

func (e *environment) IsSandbox() bool {
	return e.value.ValueOrZero() == EnvironmentSandbox
}

// IsSandBox reports whether the environment is sandboxed.
//
// Deprecated: use IsSandbox.
func (e *environment) IsSandBox() bool {
	return e.IsSandbox()
}

// IsKnown reports whether RevenueCat documents the environment as one of the Environment constants.
// Environments added by RevenueCat later (and null values) return false but remain decodable.
func (e environment) IsKnown() bool {
	return contains(validEnvironmentValues, e.value.ValueOrZero())
}

// MarshalJSON serializes an environment to JSON.
func (e environment) MarshalJSON() ([]byte, error) {
	return e.value.MarshalJSON()
}

// UnmarshalJSON deserializes an environment from JSON.
// JSON null and empty strings decode to the zero value, and undocumented values are
// preserved so that new RevenueCat values do not break decoding; see IsKnown.
func (e *environment) UnmarshalJSON(b []byte) error {
	v, err := unmarshalEnumValue(b, false)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of environment: %w", err)
	}
	e.value = v

	return nil
}
