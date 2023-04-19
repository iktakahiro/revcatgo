package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestEnvironmentUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{`"SANDBOX"`, "SANDBOX", nil},
		{`"PRODUCTION"`, "PRODUCTION", nil},
		{`"INVALID"`, "", errors.New("")},
		{`1`, "", errors.New("")},
		{`null`, "", errors.New("")},
	}

	for _, c := range cases {
		var e environment
		b := []byte(c.in)
		err := json.Unmarshal(b, &e)

		assert.Equal(t, c.expected, e.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
