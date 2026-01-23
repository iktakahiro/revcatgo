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
	// StorePlayStore identifies purchases originating from the Google Play Store.
	StorePlayStore = "PLAY_STORE"
	// StoreAppStore identifies purchases originating from the Apple App Store.
	StoreAppStore = "APP_STORE"
	// StoreStripe identifies purchases processed through Stripe.
	StoreStripe = "STRIPE"
	// StoreMacAppStore identifies purchases originating from the Mac App Store.
	StoreMacAppStore = "MAC_APP_STORE"
	// StorePromotional identifies purchases granted via promotional credits.
	StorePromotional = "PROMOTIONAL"
	// StoreRCBilling identifies purchases processed through RevenueCat Web Billing.
	StoreRCBilling = "RC_BILLING"
)

var validStoreValues = []string{
	StorePlayStore,
	StoreAppStore,
	StoreStripe,
	StoreMacAppStore,
	StorePromotional,
	StoreRCBilling,
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
