# go-sql-check

## Installation

```
go install github.com/ledgera-io/go-sql-check@latest
```

And application code:
```go
// main.go
package main

func main() {
	query := `--sql
SELECT some_nonexistent_field FROM products`
    // run query here...
}
```

We can run `go-sql-check` on it to check if all queries in the code are valid:
```bash
go-sql-check main.go
```
