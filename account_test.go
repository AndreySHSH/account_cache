package account_cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAccountAddTransactions - add transaction and check balance
func TestAccountAddTransactions(t *testing.T) {
	cases := []struct {
		Name     string
		UserName string
		Score    float64
		Expected float64
	}{
		0: {
			Name:     "Balance test_user 10x100",
			UserName: "test_user",
			Score:    10,
			Expected: 10000,
		},
		1: {
			Name:     "Balance test_user 100x1000",
			UserName: "test_user",
			Score:    100,
			Expected: 100000,
		},
		2: {
			Name:     "Balance test_user -100x1000",
			UserName: "test_user",
			Score:    -1000,
			Expected: -1000000,
		},
	}

	for _, c := range cases {
		account := Init(100000)

		for i := 0; i < 1000; i++ {
			account.Transaction(c.UserName, c.Score)
		}

		balance := account.SyncBalance(c.UserName)
		assert.Equal(t, balance, c.Expected)
	}
}
