package account_cache

import (
	"errors"
	"sync"

	"github.com/lithammer/shortuuid"
)

type User struct {
	meta any
	link *map[string]transaction
	key  int
}

type Accounts struct {
	mu           sync.RWMutex
	accounts     []User
	transactions []*map[string]transaction
}

type transaction struct {
	method string
	amount float64
}

func Init() Accounts {
	return Accounts{
		accounts:     []User{},
		transactions: []*map[string]transaction{},
	}
}

func (a *Accounts) GetAccountBalance(user any) (*User, float64, error) {
	for _, userAccount := range a.accounts {
		if userAccount.meta == user {
			var balance float64
			for id := range *userAccount.link {
				balance += (*userAccount.link)[id].amount
			}

			return &userAccount, balance, nil
		}
	}
	return nil, 0, errors.New("account not found")
}

func (a *Accounts) AddTransaction(user any, score float64) {
	tid := shortuuid.New()
	tr := transaction{
		method: "add",
		amount: score,
	}
	trm := &map[string]transaction{
		tid: tr,
	}

	for _, userAccount := range a.accounts {
		if userAccount.meta == user {
			a.mu.RLocker()
			a.mu.Lock()
			(*userAccount.link)[tid] = tr
			a.mu.Unlock()
			a.mu.RUnlock()
			return
		}
	}

	a.transactions = append(a.transactions, trm)
	for key, value := range a.transactions {
		if value == trm {
			a.accounts = append(a.accounts, User{
				meta: user,
				link: a.transactions[key],
				key:  key,
			})
			break
		}
	}
}
