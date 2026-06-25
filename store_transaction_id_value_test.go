package revcatgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreTransactionIDUnMarshal(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{`"GPA.6801-7988-0152-76034..5"`, "GPA.6801-7988-0152-76034..5"},
		{`1000000652379790`, "1000000652379790"},
		{`"a42db3af39530cb82b17eaf9c6576393"`, "a42db3af39530cb82b17eaf9c6576393"},
		{`null`, ""},
	}

	for _, c := range cases {
		var s storeTransactionID
		err := json.Unmarshal([]byte(c.in), &s)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, s.String())
	}
}
