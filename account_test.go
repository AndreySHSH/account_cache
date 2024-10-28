package account_cache

import (
	"fmt"
	"testing"
)

type Us struct {
	id   string
	name string
}

func TestAccount_Add(t *testing.T) {
	account := Init()
	account.Transaction("1", 1)

	for i := 0; i < 1000; i++ {
		//t.Run("SubTest", func(st *testing.T) {
		//	st.Parallel()
		account.Transaction("123", 1)
		//})
	}

	ac := account.SyncBalance("123")
	fmt.Println(ac)
}
