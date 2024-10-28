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

type transaction struct {
	method string
	amount float64
}

type Tr struct {
	User  any
	Score float64
}

type Engin struct {
	Queue chan Tr

	accounts     []User
	transactions []*sync.Map
}

func Init() *Engin {
	e := &Engin{
		Queue: make(chan Tr, 100000),

		accounts:     []User{},
		transactions: []*sync.Map{},
	}

	go e.readQueue()
	return e
}

func (a *Engin) GetAccountBalance(user any) (*User, float64, error) {
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

func (a *Engin) AddTransaction(user any, score float64) {
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

func (a *Engin) readQueue() {
	for {
		select {
		case tr := <-a.Queue:
			a.AddTransaction(tr.User, tr.Score)
		}
	}
}
