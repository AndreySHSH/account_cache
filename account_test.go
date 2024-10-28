package account_cache

import (
	"fmt"
	"testing"
	"time"
)

type Us struct {
	id   string
	name string
}

func TestAccount_Add(t *testing.T) {
	//t.Parallel()
	account := Init()
	account.Transaction("1", 1)

	for i := 0; i < 1000; i++ {
		go account.Transaction("123", 1)
	}

	ac, f, err := account.Balance("123")
	fmt.Println(ac, f, err)
	time.Sleep(1 * time.Second)

}
