package account_cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAccountAddTransactions - add transaction and check balance
func TestAccountAddTransactions(t *testing.T) {
	t.Parallel()

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

		if wallet := account.CreateWallet(c.UserName); wallet == nil {
			account.Transaction(c.UserName, 0)
		}

		for i := 0; i < 1000; i++ {
			account.Transaction(c.UserName, c.Score)
		}

		balance := account.SyncBalance(c.UserName)

		assert.Equal(t, balance, c.Expected)
	}
}

// TestEngin_SetCurrentAccountBalance - test on SetCurrentAccountBalance work
func TestEngin_SetCurrentAccountBalance(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name       string
		UserName   string
		Score      float64
		Expected   float64
		SetCurrent float64
	}{
		0: {
			Name:       "Balance test_user 10x100",
			UserName:   "test_user",
			Score:      10,
			Expected:   10,
			SetCurrent: 10,
		},
		1: {
			Name:       "Balance test_user 100x1000",
			UserName:   "test_user",
			Score:      100,
			Expected:   100,
			SetCurrent: 100,
		},
		2: {
			Name:       "Balance test_user -100x1000",
			UserName:   "test_user",
			Score:      100,
			Expected:   20000,
			SetCurrent: 20000,
		},
	}

	for _, c := range cases {
		account := Init(100000)

		if wallet := account.CreateWallet(c.UserName); wallet == nil {
			account.Transaction(c.UserName, 0)
		}

		for i := 0; i < 1000; i++ {
			account.Transaction(c.UserName, c.Score)
		}

		account.SetCurrentAccountBalance(c.UserName, c.SetCurrent)

		balance := account.SyncBalance(c.UserName)

		assert.Equal(t, balance, c.Expected)
	}
}
