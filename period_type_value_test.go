package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPeriodType(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"TRIAL", "TRIAL", nil},
		{"NORMAL", "NORMAL", nil},
		{"PREPAID", "PREPAID", nil},
		{"INVALID", "", errors.New("periodType value should be one of the following: TRIAL,INTRO,NORMAL,PROMOTIONAL,PREPAID, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newPeriodType(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestPeriodTypeUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		known    bool
		wantErr  bool
	}{
		{`"TRIAL"`, "TRIAL", true, false},
		{`"NORMAL"`, "NORMAL", true, false},
		{`"prepaid"`, "PREPAID", true, false},
		{`"FUTURE_TYPE"`, "FUTURE_TYPE", false, false},
		{`"future_type"`, "FUTURE_TYPE", false, false},
		{`""`, "", false, false},
		{`null`, "", false, false},
		{`1`, "", false, true},
	}

	for _, c := range cases {
		var p periodType
		err := json.Unmarshal([]byte(c.in), &p)
		if c.wantErr {
			require.Error(t, err, c.in)
			continue
		}
		require.NoError(t, err, c.in)
		assert.Equal(t, c.expected, p.String(), c.in)
		assert.Equal(t, c.known, p.IsKnown(), c.in)
	}
}

func TestPeriodTypeIsKnown(t *testing.T) {
	known, err := newPeriodType(PeriodTypeTrial)
	require.NoError(t, err)
	assert.True(t, known.IsKnown())

	var unknown periodType
	require.NoError(t, json.Unmarshal([]byte(`"FUTURE_TYPE"`), &unknown))
	assert.False(t, unknown.IsKnown())
	assert.Equal(t, "FUTURE_TYPE", unknown.String())

	var zero periodType
	assert.False(t, zero.IsKnown())
}

func TestPeriodTypeMarshalRoundTrip(t *testing.T) {
	for _, in := range []string{
		`"TRIAL"`,
		`"FUTURE_TYPE"`,
		`null`,
	} {
		var p periodType
		require.NoError(t, json.Unmarshal([]byte(in), &p))

		out, err := json.Marshal(p)
		require.NoError(t, err)
		assert.JSONEq(t, in, string(out))
	}
}
