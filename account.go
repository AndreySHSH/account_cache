// Package account_cache - synchronous and asynchronous work with account balances
package account_cache

import (
	"sync"

	"github.com/lithammer/shortuuid"
)

// Engin - engine for caching user account balance
type Engin struct {
	queue        chan task
	read         bool
	accounts     []User
	transactions []*sync.Map
}

// Init - Engin struct
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

// Transaction - public method for adding transactions
func (engin *Engin) Transaction(user any, score float64) {
	engin.queue <- task{
		User:         user,
		Score:        score,
		Notification: nil,
	}
}

// AsyncBalance - asynchronous method for obtaining account balance
func (engin *Engin) AsyncBalance(user any) float64 {
	for _, userAccount := range engin.accounts {
		if userAccount.meta == user {
			var balance float64
			userAccount.link.Range(func(key, value any) bool {
				balance = balance + value.(transaction).amount
				return true
			})
			return balance
		}
	}
	return 0
}

// SyncBalance - synchronous method for obtaining account balance
func (engin *Engin) SyncBalance(user any) float64 {
	notification := make(chan float64)

	engin.queue <- task{
		User:         user,
		Score:        0,
		Notification: notification,
	}
	return <-notification
}

// add - adding a transaction to storage
func (engin *Engin) add(user any, score float64) {
	transactionId := shortuuid.New()
	userTransaction := transaction{
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
	transactionStore := engin.transactions[len(engin.transactions)-1]
	if _, ok := transactionStore.Load(transactionId); ok {
		engin.accounts = append(engin.accounts, User{
			meta: user,
			link: transactionStore,
		})
	}
}

// worker - queue worker
func (engin *Engin) worker() {
	for {
		select {
		case tr := <-engin.queue:
			if tr.Notification != nil {
				tr.Notification <- engin.AsyncBalance(tr.User)
				continue
			}
			engin.add(tr.User, tr.Score)
		}
	}
}
