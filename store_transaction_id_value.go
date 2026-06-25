package revcatgo

import (
	"bytes"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

// storeTransactionID represents the store_transaction_id field returned by the
// RevenueCat REST API. The App Store reports it as a JSON number while the
// other stores report it as a JSON string, so the value is normalized to a
// string.
type storeTransactionID struct {
	value null.String
}

// String returns the store transaction id.
func (s *storeTransactionID) String() string {
	return s.value.ValueOrZero()
}

func (s storeTransactionID) MarshalJSON() ([]byte, error) {
	return s.value.MarshalJSON()
}

// UnmarshalJSON deserializes a store transaction id from either a JSON string
// or a JSON number.
func (s *storeTransactionID) UnmarshalJSON(b []byte) error {
	// Numbers are reported unquoted (e.g. by the App Store); quote them so the
	// raw token is preserved exactly and parsed as a string.
	if len(b) > 0 && b[0] != '"' && !bytes.Equal(b, []byte("null")) {
		b = append([]byte{'"'}, append(b, '"')...)
	}
	if err := s.value.UnmarshalJSON(b); err != nil {
		return fmt.Errorf("failed to unmarshal the value of store transaction id: %w", err)
	}

	return nil
}
