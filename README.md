# account_cache

#### Package account_cache - synchronous work with account balances

Use SyncBalance
```go

bufferSize := 100000

user := "test"
score := 1000

account := Init(bufferSize) // buffer queue size.

if wallet := account.CreateWallet(user); wallet == nil {
    // if wallet == nil your wallet created now. 
    // wallet ready to create transaction.
    account.Transaction(user, 100) // add first default balance.
}

for i := 0; i < 1000; i++ {
    account.Transaction(user, score)
}

balance := account.SyncBalance(user)
```

Benchmark
```text
goos: darwin
goarch: arm64
pkg: github.com/AndreySHSH/account_cache
cpu: Apple M2
BenchmarkGetBalance-8              10000            117022 ns/op          85.45 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117594 ns/op          85.04 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117286 ns/op          85.26 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117219 ns/op          85.31 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117185 ns/op          85.34 MB/s        2970 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            118355 ns/op          84.49 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117808 ns/op          84.88 MB/s        2972 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117074 ns/op          85.42 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            116800 ns/op          85.62 MB/s        2971 B/op        119 allocs/op
BenchmarkGetBalance-8              10000            117076 ns/op          85.41 MB/s        2970 B/op        119 allocs/op
```