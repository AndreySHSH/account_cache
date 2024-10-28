package account_cache

import (
	"errors"
	"sync"

	"github.com/lithammer/shortuuid"
)

type Engin struct {
	queue        chan task
	read         bool
	accounts     []User
	transactions []*sync.Map
}

func Init() *Engin {
	e := &Engin{
		queue: make(chan task, 100000),

		accounts:     []User{},
		transactions: []*sync.Map{},
	}

	go e.worker()

	return e
}

func (engin *Engin) Transaction(user any, score float64) {
	engin.queue <- task{
		User:  user,
		Score: score,
	}
}

func (engin *Engin) Balance(user any) (*User, float64, error) {
	for {
		if engin.read {
			for _, userAccount := range engin.accounts {
				if userAccount.meta == user {
					var balance float64
					userAccount.link.Range(func(key, value any) bool {
						balance = balance + value.(transaction).amount
						return true
					})

					return &userAccount, balance, nil
				}
			}
			return nil, 0, errors.New("account not found")
		}
	}
}

func (engin *Engin) add(user any, score float64) {
	transactionId := shortuuid.New()
	userTransaction := transaction{
		method: "add",
		amount: score,
	}

	for _, userAccount := range engin.accounts {
		if userAccount.meta == user {
			userAccount.link.Store(transactionId, userTransaction)
			return
		}
	}

	data := sync.Map{}
	data.Store(transactionId, userTransaction)

	engin.transactions = append(engin.transactions, &data)
	for keyUserTransactionStore, transactionStore := range engin.transactions {
		if _, ok := transactionStore.Load(transactionId); ok {
			engin.accounts = append(engin.accounts, User{
				meta: user,
				link: transactionStore,
				key:  keyUserTransactionStore,
			})
			break
		}
	}
}

func (engin *Engin) lock() {
	engin.read = false
}

func (engin *Engin) unlock() {
	engin.read = true
}

func (engin *Engin) worker() {
	for {
		select {
		case tr := <-engin.queue:
			engin.lock()
			engin.add(tr.User, tr.Score)
			engin.unlock()
		}
	}
}
