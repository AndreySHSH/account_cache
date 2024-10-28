package account_cache

import (
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
		User:         user,
		Score:        score,
		Notification: nil,
	}
}

func (engin *Engin) AsyncBalance(user any) float64 {
	for {
		if engin.read {
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
	}
}

func (engin *Engin) SyncBalance(user any) float64 {
	notification := make(chan float64)

	engin.queue <- task{
		User:         user,
		Score:        0,
		Notification: notification,
	}
	return <-notification
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
			if tr.Notification != nil {
				tr.Notification <- engin.AsyncBalance(tr.User)
				continue
			}
			engin.lock()
			engin.add(tr.User, tr.Score)
			engin.unlock()
		}
	}
}
