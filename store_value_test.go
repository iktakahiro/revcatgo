package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"PLAY_STORE", "PLAY_STORE", nil},
		{"APP_STORE", "APP_STORE", nil},
		{"AMAZON", "AMAZON", nil},
		{"PADDLE", "PADDLE", nil},
		{"RC_BILLING", "RC_BILLING", nil},
		{"ROKU", "ROKU", nil},
		{"TEST_STORE", "TEST_STORE", nil},
		{"INVALID", "", errors.New("store value should be one of the following: AMAZON, APP_STORE, MAC_APP_STORE, PADDLE, PLAY_STORE, RC_BILLING, ROKU, PROMOTIONAL, STRIPE, TEST_STORE, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newStore(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestStoreUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		known    bool
		wantErr  bool
	}{
		{`"PLAY_STORE"`, "PLAY_STORE", true, false},
		{`"APP_STORE"`, "APP_STORE", true, false},
		{`"amazon"`, "AMAZON", true, false},
		{`"paddle"`, "PADDLE", true, false},
		{`"rc_billing"`, "RC_BILLING", true, false},
		{`"GALAXY"`, "GALAXY", false, false},
		{`"galaxy"`, "GALAXY", false, false},
		{`""`, "", false, false},
		{`null`, "", false, false},
		{`1`, "", false, true},
	}

	for _, c := range cases {
		var s store
		err := json.Unmarshal([]byte(c.in), &s)
		if c.wantErr {
			require.Error(t, err, c.in)
			continue
		}
		require.NoError(t, err, c.in)
		assert.Equal(t, c.expected, s.String(), c.in)
		assert.Equal(t, c.known, s.IsKnown(), c.in)
	}
}

func TestStoreIsKnown(t *testing.T) {
	known, err := newStore(StoreAppStore)
	require.NoError(t, err)
	assert.True(t, known.IsKnown())

	var unknown store
	require.NoError(t, json.Unmarshal([]byte(`"GALAXY"`), &unknown))
	assert.False(t, unknown.IsKnown())
	assert.Equal(t, "GALAXY", unknown.String())

	var zero store
	assert.False(t, zero.IsKnown())
}

func TestStoreMarshalRoundTrip(t *testing.T) {
	for _, in := range []string{
		`"APP_STORE"`,
		`"GALAXY"`,
		`null`,
	} {
		var s store
		require.NoError(t, json.Unmarshal([]byte(in), &s))

		out, err := json.Marshal(s)
		require.NoError(t, err)
		assert.JSONEq(t, in, string(out))
	}
}
