package revcatgo

import (
	"errors"
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
// Deprecated: use IsSandbox.
func (e *environment) IsSandBox() bool {
	return e.IsSandbox()
}

// MarshalJSON serializes an environment to JSON.
func (e environment) MarshalJSON() ([]byte, error) {
	return e.value.MarshalJSON()
}

// UnmarshalJSON deserializes an environment from JSON.
func (e *environment) UnmarshalJSON(b []byte) error {
	v := &environment{}
	err := v.value.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of environment: %w", err)
	}
	if !v.value.Valid {
		return errors.New("environment is a required field")
	}
	_e, err := newEnvironment(v.value.ValueOrZero())
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of environment: %w", err)
	}
	e.value = _e.value

	return nil
}
