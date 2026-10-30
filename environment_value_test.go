package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvironment(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"SANDBOX", "SANDBOX", nil},
		{"PRODUCTION", "PRODUCTION", nil},
		{"INVALID", "", errors.New("environment value should be one of the following: SANDBOX,PRODUCTION, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newEnvironment(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestEnvironmentUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		known    bool
		wantErr  bool
	}{
		{`"SANDBOX"`, "SANDBOX", true, false},
		{`"PRODUCTION"`, "PRODUCTION", true, false},
		{`"STAGING"`, "STAGING", false, false},
		{`"sandbox"`, "sandbox", false, false},
		{`""`, "", false, false},
		{`null`, "", false, false},
		{`1`, "", false, true},
	}

	for _, c := range cases {
		var e environment
		err := json.Unmarshal([]byte(c.in), &e)
		if c.wantErr {
			require.Error(t, err, c.in)
			continue
		}
		require.NoError(t, err, c.in)
		assert.Equal(t, c.expected, e.String(), c.in)
		assert.Equal(t, c.known, e.IsKnown(), c.in)
	}
}

func TestEnvironmentIsKnown(t *testing.T) {
	known, err := newEnvironment(EnvironmentProduction)
	require.NoError(t, err)
	assert.True(t, known.IsKnown())

	var unknown environment
	require.NoError(t, json.Unmarshal([]byte(`"STAGING"`), &unknown))
	assert.False(t, unknown.IsKnown())
	assert.Equal(t, "STAGING", unknown.String())

	var zero environment
	assert.False(t, zero.IsKnown())
}

func TestEnvironmentMarshalRoundTrip(t *testing.T) {
	for _, in := range []string{
		`"PRODUCTION"`,
		`"STAGING"`,
		`null`,
	} {
		var e environment
		require.NoError(t, json.Unmarshal([]byte(in), &e))

		out, err := json.Marshal(e)
		require.NoError(t, err)
		assert.JSONEq(t, in, string(out))
	}
}

func TestEnvironmentHelpers(t *testing.T) {
	sandbox, err := newEnvironment(EnvironmentSandbox)
	require.NoError(t, err)
	assert.True(t, sandbox.IsSandbox())
	assert.False(t, sandbox.IsProduction())

	production, err := newEnvironment(EnvironmentProduction)
	require.NoError(t, err)
	assert.True(t, production.IsProduction())
	assert.False(t, production.IsSandbox())
}

func TestEnvironmentIsSandBox(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  bool
	}{
		{"sandbox", `"SANDBOX"`, true},
		{"production", `"PRODUCTION"`, false},
		{"null", `null`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var value environment
			require.NoError(t, json.Unmarshal([]byte(tc.input), &value))
			assert.Equal(t, tc.want, value.IsSandBox())
		})
	}
}
