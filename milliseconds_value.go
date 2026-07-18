package revcatgo

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

type milliseconds struct {
	value null.Int
}

const (
	// 1973-03-03T09:46:40Z in milliseconds since the Unix epoch.
	millisecondsThreshold = 100_000_000_000
)

func newMilliseconds(v null.Int) (*milliseconds, error) {
	if v.ValueOrZero() != 0 && v.ValueOrZero() < millisecondsThreshold {
		return &milliseconds{}, fmt.Errorf("milliseconds should be greater than %v", millisecondsThreshold)
	}
	return &milliseconds{value: v}, nil
}

func (m *milliseconds) Int64() int64 {
	return m.value.ValueOrZero()
}

func (m *milliseconds) NullInt() null.Int {
	return m.value
}

func (m *milliseconds) String() string {
	return fmt.Sprint(m.Int64())
}

func (m *milliseconds) DateTime() time.Time {
	return time.UnixMilli(m.value.ValueOrZero())
}

func (m milliseconds) MarshalJSON() ([]byte, error) {
	return m.value.MarshalJSON()
}

func (m *milliseconds) UnmarshalJSON(b []byte) error {
	v := &milliseconds{}
	err := v.value.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of milliseconds: %w", err)
	}
	_m, err := newMilliseconds(v.value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of milliseconds: %w", err)
	}
	m.value = _m.value

	return nil
}
