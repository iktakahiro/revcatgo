package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPeriodType(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"TRIAL", "TRIAL", nil},
		{"NORMAL", "NORMAL", nil},
		{"INVALID", "", errors.New("periodType value should be one of the following: TRIAL,INTRO,NORMAL,PROMOTIONAL, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newPeriodType(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestPeriodTypeUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{`"TRIAL"`, "TRIAL", nil},
		{`"NORMAL"`, "NORMAL", nil},
		{`"INVALID"`, "", errors.New("")},
		{`1`, "", errors.New("")},
		{`null`, "", errors.New("")},
	}

	for _, c := range cases {
		var p periodType
		b := []byte(c.in)
		err := json.Unmarshal(b, &p)

		if err == nil {
			assert.Equal(t, c.expected, p.String())
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
