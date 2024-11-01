package account_cache

import (
	"testing"
)

func BenchmarkGetBalance(b *testing.B) {
	account := Init(100000)

	b.SetBytes(int64(10000))
	b.ResetTimer()

	if wallet := account.CreateWallet("test_user"); wallet == nil {
		account.Transaction("test_user", 0)
	}
	
	for i := 0; i < b.N; i++ {
		account.Transaction("test_user", float64(i))
	}

	for i := 0; i < b.N; i++ {
		_ = account.SyncBalance("test_user")
	}
}
