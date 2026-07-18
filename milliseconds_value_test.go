package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestNewMilliSeconds(t *testing.T) {
	cases := []struct {
		in       int64
		expected time.Time
		err      error
	}{
		{1605526568123, time.UnixMilli(1605526568123), nil},
		{1605526568, time.UnixMilli(0), errors.New("milliseconds should be greater than 100000000000")},
		{0, time.Unix(0, 0), nil},
	}

	for _, c := range cases {
		actual, err := newMilliseconds(null.IntFrom(c.in))
		if err == nil {
			assert.Equal(t, c.expected, actual.DateTime())
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestMilliSecondsUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected int64
		err      error
	}{
		{`1605526568123`, 1605526568123, nil},
		{`1605526568`, 0, errors.New("failed to unmarshal the value of milliseconds: milliseconds should be greater than 100000000000")},
		{`null`, 0, nil},
	}

	for _, c := range cases {
		var m milliseconds
		b := []byte(c.in)
		err := json.Unmarshal(b, &m)

		if err == nil {
			assert.Equal(t, c.expected, m.Int64())
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}
