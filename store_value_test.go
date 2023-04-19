package revcatgo

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{"PLAY_STORE", "PLAY_STORE", nil},
		{"APP_STORE", "APP_STORE", nil},
		{"INVALID", "", errors.New("store value should be one of the following: PLAY_STORE, APP_STORE, STRIPE, MAC_APP_STORE, PROMOTIONAL, got INVALID")},
	}

	for _, c := range cases {
		actual, err := newStore(c.in)
		assert.Equal(t, c.expected, actual.String())
		if c.err == nil {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.err.Error())
		}
	}
}

func TestStoreUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
		err      error
	}{
		{`"PLAY_STORE"`, "PLAY_STORE", nil},
		{`"APP_STORE"`, "APP_STORE", nil},
		{`"INVALID"`, "", errors.New("")},
		{`1`, "", errors.New("")},
		{`null`, "", errors.New("")},
	}

	for _, c := range cases {
		var s store
		b := []byte(c.in)
		err := json.Unmarshal(b, &s)

		if err == nil {
			assert.Equal(t, c.expected, s.String())
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
