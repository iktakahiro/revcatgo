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
		{1605526568000, time.Unix(1605526568000/1000, 0), nil},
		{16055265, time.Unix(0, 0), errors.New("milliseconds should be grater than 100000000")},
		{0, time.Unix(0, 0), nil},
	}

	for _, c := range cases {
		actual, err := newMilliseconds(null.IntFrom(c.in))
		if err == nil {
			assert.Equal(t, c.expected, actual.DateTime())
			assert.Nil(t, err)
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
		{`1605526568000`, 1605526568000, nil},
		{`16055265`, 0, errors.New("failed to unmarshal the value of milliseconds: milliseconds should be grater than 100000000")},
		{`null`, 0, nil},
	}

	for _, c := range cases {
		var m milliseconds
		b := []byte(c.in)
		err := json.Unmarshal(b, &m)

		if err == nil {
			assert.Equal(t, c.expected, m.Int64())
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}
