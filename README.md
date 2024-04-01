[![GoDoc](https://pkg.go.dev/badge/github.com/blizzy78/gopii)](https://pkg.go.dev/github.com/blizzy78/gopii)


gopii
=====

A Go package that provides functions to sanitize personally identifiable information (PII) from text.

```go
import "github.com/blizzy78/gopii"
```


Code example
------------

```go
// This function does the actual work.
// In this example, we're not using the Context, but you really always should.
worker := func(_ context.Context) {
	time.Sleep(100 * time.Millisecond)
}

goroutines := New()

// Start a new goroutine.
goroutines.Go(context.Background(), worker)

// Cancel all goroutines' contexts, and wait for them to finish.
_ = goroutines.CancelAll(context.Background(), true)
```


License
-------

This package is licensed under the MIT license.
