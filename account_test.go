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
	user := Us{
		id:   "1",
		name: "andrey",
	}

	account := Init()

	for i := 0; i < 1000; i++ {
		account.AddTransaction(user, 200)

	}
	//account.AddTransaction(user, 200)
	//account.AddTransaction(user, 200)
	//account.AddTransaction(user, 200)
	//account.AddTransaction(user, -200)

	ac, f, err := account.GetAccountBalance(user)
	fmt.Println(ac, f, err)

}
