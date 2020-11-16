package revcatgo

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type store struct {
	value null.String
}

const (
	storePlayStore   = "PLAY_STORE"
	storeAppStore    = "APP_STORE"
	storeStripe      = "STRIPE"
	storeMacAppStore = "MAC_APP_STORE"
	storePromotional = "PROMOTIONAL"
)

var validStoreValues = []string{
	storePlayStore,
	storeAppStore,
	storeStripe,
	storeMacAppStore,
	storePromotional,
}

func newStore(s string) (*store, error) {
	if !contains(validStoreValues, s) {
		return &store{}, errors.New("store value should be one of the following: " + strings.Join(validStoreValues, ", "))
	}
	return &store{value: null.StringFrom(s)}, nil
}

func (s *store) String() string {
	return s.value.ValueOrZero()
}

func (s store) MarshalJSON() ([]byte, error) {
	return s.value.MarshalJSON()
}

// UnmarshalJSON deserializes a store from JSON
func (s *store) UnmarshalJSON(b []byte) error {
	v := &store{}
	err := v.value.UnmarshalJSON(b)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of store: %w", err)
	}
	if !v.value.Valid {
		return errors.New("store is a required field")
	}
	_s, err := newStore(v.value.ValueOrZero())
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of store: %w", err)
	}
	s.value = _s.value

	return nil
}
