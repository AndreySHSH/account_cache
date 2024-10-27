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

	for i := 0; i < 1000; i++ {
		account.AddTransaction("123", 200)
	}
	//account.AddTransaction(user, 200)
	//account.AddTransaction(user, 200)
	//account.AddTransaction(user, 200)
	//account.AddTransaction(user, -200)

	ac, f, err := account.GetAccountBalance("123")
	fmt.Println(ac, f, err)

}
