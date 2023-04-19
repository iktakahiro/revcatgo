package revcatgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestNewPrice(t *testing.T) {
	cases := []struct {
		in       float64
		expected float64
		err      error
	}{
		{500, 500, nil},
		{0, 0, nil},
	}

	for _, c := range cases {
		actual, err := newPrice(null.FloatFrom(c.in))
		assert.Equal(t, c.expected, actual.Float64())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
	cases2 := []struct {
		in       float64
		expected bool
		err      error
	}{
		{500, false, nil},
		{0, true, nil},
	}

	for _, c := range cases2 {
		actual, err := newPrice(null.FloatFrom(c.in))
		if err == nil {
			assert.Equal(t, c.expected, actual.IsFreeTrial())
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestPriceUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected float64
		err      error
	}{
		{`500`, 500, nil},
		{`0`, 0, nil},
		{`null`, 0, nil},
	}

	for _, c := range cases {
		var p price
		b := []byte(c.in)
		err := json.Unmarshal(b, &p)

		assert.Equal(t, c.expected, p.Float64())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
