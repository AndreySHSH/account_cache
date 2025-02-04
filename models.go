package account_cache

import (
	"sync"
)

type User struct {
	meta any
	link *sync.Map
}

func (user *User) HistoryLength() (length int64) {
	user.link.Range(func(key, value interface{}) bool {
		length += 1

		return true
	})

	return
}

type transaction struct {
	amount float64
}

type task struct {
	user         any
	score        float64
	notification chan float64
	transaction  chan string
	tt           string
}
