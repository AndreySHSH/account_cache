package account_cache

import "sync"

type User struct {
	meta any
	link *sync.Map
}

type transaction struct {
	amount float64
}

type task struct {
	user         any
	score        float64
	notification chan float64
	transaction  chan string
}
