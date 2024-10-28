# account_cache

#### Package account_cache - synchronous and asynchronous work with account balances

Use SyncBalance
```go

bufferSize := 100000

user := "test"
score := 1000

account := Init(bufferSize) // buffer queue size

for i := 0; i < 1000; i++ {
    account.Transaction(user, score)
}

balance := account.SyncBalance(user)
```

Use AsyncBalance
```go
bufferSize := 100000

user := "test"
score := 1000

account := Init(bufferSize) // buffer queue size

// use if the user account has not yet been created
if account.AsyncBalance(user) == 0 { 
    account.Transaction(user, score)
}

for i := 0; i < 1000; i++ {
    account.Transaction(user, 2000) // add new transaction
}

balance := account.SyncBalance(user)
```