package account_cache

import "sync"

type User struct {
	meta any
	link *sync.Map
	key  int
}

type transaction struct {
	method string
	amount float64
}

type task struct {
	User         any
	Score        float64
	Notification chan float64
}
