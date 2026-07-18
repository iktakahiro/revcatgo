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
	// StoreAmazon identifies purchases originating from the Amazon Appstore.
	StoreAmazon = "AMAZON"
	// StoreAppStore identifies purchases originating from the Apple App Store.
	StoreAppStore = "APP_STORE"
	// StoreMacAppStore identifies purchases originating from the Mac App Store.
	StoreMacAppStore = "MAC_APP_STORE"
	// StorePaddle identifies purchases originating from the Paddle Store.
	StorePaddle = "PADDLE"
	// StorePlayStore identifies purchases originating from the Google Play Store.
	StorePlayStore = "PLAY_STORE"
	// StoreRCBilling identifies purchases granted via the Revenue Cat web Store.
	StoreRCBilling = "RC_BILLING"
	// StoreRoku identifies purchases granted via the Roku Store.
	StoreRoku = "ROKU"
	// StorePromotional identifies purchases granted via promotional credits.
	StorePromotional = "PROMOTIONAL"
	// StoreStripe identifies purchases processed through Stripe.
	StoreStripe = "STRIPE"
	// StoreTestStore identifies purchases via the test store
	StoreTestStore = "TEST_STORE"
)

var validStoreValues = []string{
	StoreAmazon,
	StoreAppStore,
	StoreMacAppStore,
	StorePaddle,
	StorePlayStore,
	StoreRCBilling,
	StoreRoku,
	StorePromotional,
	StoreStripe,
	StoreTestStore,
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
