package account_cache

import (
	"errors"
	"sync"

	"github.com/lithammer/shortuuid"
)

type User struct {
	meta any
	link *sync.Map
	key  int
}

type Accounts struct {
	mu           sync.RWMutex
	accounts     []User
	transactions []*sync.Map
}

type transaction struct {
	method string
	amount float64
}

func Init() Accounts {
	return Accounts{
		accounts:     []User{},
		transactions: []*sync.Map{},
	}
}

func (a *Accounts) GetAccountBalance(user any) (*User, float64, error) {
	for _, userAccount := range a.accounts {
		if userAccount.meta == user {
			var balance float64
			userAccount.link.Range(func(key, value any) bool {
				balance += value.(transaction).amount
				return true
			})

			return &userAccount, balance, nil
		}
	}
	return nil, 0, errors.New("account not found")
}

func (a *Accounts) AddTransaction(user any, score float64) {
	transactionId := shortuuid.New()
	userTransaction := transaction{
		method: "add",
		amount: score,
	}

	for _, userAccount := range a.accounts {
		if userAccount.meta == user {
			userAccount.link.Store(transactionId, userTransaction)
			return
		}
	}

	data := sync.Map{}
	data.Store(transactionId, userTransaction)

	a.transactions = append(a.transactions, &data)
	for keyUserTransactionStore, transactionStore := range a.transactions {

		if _, ok := transactionStore.Load(transactionId); ok {
			a.accounts = append(a.accounts, User{
				meta: user,
				link: transactionStore,
				key:  keyUserTransactionStore,
			})
			break
		}
	}
}
