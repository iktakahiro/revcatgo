package revcatgo

import (
	"slices"
	"strings"

	"gopkg.in/guregu/null.v4"
)

func contains(s []string, e string) bool {
	return slices.Contains(s, e)
}

// unmarshalEnumValue decodes a JSON string or null for an enum-like value type.
// JSON null and empty strings yield the zero value, and any other string is kept
// as-is (upper-cased when upper is true) so that values RevenueCat introduces
// later do not break decoding. Non-string JSON values return an error.
func unmarshalEnumValue(b []byte, upper bool) (null.String, error) {
	var v null.String
	if err := v.UnmarshalJSON(b); err != nil {
		return null.String{}, err
	}
	if v.ValueOrZero() == "" {
		return null.String{}, nil
	}
	if upper {
		return null.StringFrom(strings.ToUpper(v.String)), nil
	}

	return v, nil
}
