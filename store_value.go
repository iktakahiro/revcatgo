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
	StoreAmazon      = "AMAZON"
	StorePlayStore   = "PLAY_STORE"
	StoreAppStore    = "APP_STORE"
	StoreStripe      = "STRIPE"
	StoreMacAppStore = "MAC_APP_STORE"
	StorePromotional = "PROMOTIONAL"
)

var validStoreValues = []string{
	StorePlayStore,
	StoreAppStore,
	StoreStripe,
	StoreMacAppStore,
	StorePromotional,
}

func newStore(s string) (*store, error) {
	if !contains(validStoreValues, s) {
		return &store{}, fmt.Errorf("store value should be one of the following: %v, got %v", strings.Join(validStoreValues, ", "), s)
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
	_s, err := newStore(strings.ToUpper(v.value.ValueOrZero()))
	if err != nil {
		return fmt.Errorf("failed to unmarshal the value of store: %w", err)
	}
	s.value = _s.value

	return nil
}
