// Package account_cache - synchronous work with account balances.
package account_cache

import (
	"errors"
	"sync"

	"github.com/lithammer/shortuuid"
)

// Engin - engine for caching user account balance.
type Engin struct {
	queue        chan task
	accounts     []User
	transactions []*sync.Map
}

// Init - Engin struct.
func Init(queueBufferSize int64) *Engin {
	engin := &Engin{
		queue:        make(chan task, queueBufferSize),
		accounts:     []User{},
		transactions: []*sync.Map{},
	}

	// run queue worker reader
	go engin.worker()

	return engin
}

// Transaction - public method for adding transactions.
func (engin *Engin) Transaction(user any, score float64) (transactionId string) {
	tid := make(chan string)
	engin.queue <- task{
		user:         user,
		score:        score,
		notification: nil,
		transaction:  tid,
	}

	transactionId = <-tid

	return
}

// SetCurrentAccountBalance - set default balance.
func (engin *Engin) SetCurrentAccountBalance(user any, score float64) (transactionId string) {
	tid := make(chan string)
	engin.queue <- task{
		user:         user,
		score:        score,
		notification: nil,
		transaction:  tid,
		action:       "set_current_balance",
	}

	transactionId = <-tid

	return
}

func (engin *Engin) CreateWallet(user any) *User {
	for _, account := range engin.accounts {
		if account.meta == user {
			return &account
		}
	}

	data := sync.Map{}

	engin.transactions = append(engin.transactions, &data)
	transactionStore := engin.transactions[len(engin.transactions)-1]

	engin.accounts = append(engin.accounts, User{
		meta: user,
		link: transactionStore,
	})

	return nil
}

// Rollback - cancel transaction.
func (engin *Engin) Rollback(user any, transactionId string) error {
	for _, userAccount := range engin.accounts {
		if userAccount.meta == user {
			userAccount.link.Delete(transactionId)

			return nil
		}
	}

	return errors.New("user not found")
}

// SyncBalance - synchronous method for obtaining account balance.
func (engin *Engin) SyncBalance(user any) float64 {
	notification := make(chan float64)

	engin.queue <- task{
		user:         user,
		score:        0,
		notification: notification,
	}

	return <-notification
}

// add - adding a transaction to storage.
func (engin *Engin) add(user any, score float64) string {
	transactionId := shortuuid.New()
	userTransaction := transaction{
		amount: score,
	}

	for _, userAccount := range engin.accounts {
		if userAccount.meta == user {
			userAccount.link.Store(transactionId, userTransaction)

			return transactionId
		}
	}

	data := sync.Map{}
	data.Store(transactionId, userTransaction)

	engin.transactions = append(engin.transactions, &data)
	transactionStore := engin.transactions[len(engin.transactions)-1]

	if _, ok := transactionStore.Load(transactionId); ok {
		engin.accounts = append(engin.accounts, User{
			meta: user,
			link: transactionStore,
		})
	}

	return transactionId
}

// current - current balance.
func (engin *Engin) current(user any, score float64) string {
	transactionId := shortuuid.New()
	userTransaction := transaction{
		amount: score,
	}
	for _, userAccount := range engin.accounts {
		if userAccount.meta == user {
			userAccount.link.Range(func(key, value interface{}) bool {
				userAccount.link.Delete(key)
				return true
			})
			userAccount.link.Store(transactionId, userTransaction)
		}
	}

	return transactionId
}

// asyncBalance - asynchronous method for obtaining account balance.
func (engin *Engin) asyncBalance(user any) float64 {
	for _, userAccount := range engin.accounts {
		if userAccount.meta == user {
			var balance float64
			userAccount.link.Range(func(key, value any) bool {
				balance += value.(transaction).amount

				return true
			})

			return balance
		}
	}

	return 0
}

// worker - queue worker.
func (engin *Engin) worker() {
	for {
		t := <-engin.queue
		if t.action == "set_current_balance" {
			t.transaction <- engin.current(t.user, t.score)

			continue
		}
		if t.notification != nil {
			t.notification <- engin.asyncBalance(t.user)

			continue
		}
		t.transaction <- engin.add(t.user, t.score)
	}
}
