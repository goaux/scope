# scope

`scope` is a Go module that provides iterators that can be used to automatically close, unlock, or cancel resources like files when they are no longer needed.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/scope.svg)](https://pkg.go.dev/github.com/goaux/scope)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/scope)](https://goreportcard.com/report/github.com/goaux/scope)

The iterator returned by this package will execute the loop body exactly once.

The defer idiom is commonly used in Go to allocate and release resources.

Deferred actions are executed when the function returns, so the release of resources is delayed until the function returns.

Sometimes that's what you want, and sometimes it's not.

When resources need to be released immediately, function literals are often used to specify the scope.

```go
func() {
    file, err := os.Open(name)
    if err != nil {
      return
    }
    defer file.Close()
    _ = file // use file here
}()
```

This package provides another solution for such situations.
The following code snippet achieves the same result as the code above.

```go
import "github.com/goaux/iter/scopeos"

for file, err := range scope.Use2(os.Open(name)) {
    // If err==nil, the file will be closed at the end of the loop body regardless of break.
    if err != nil {
        break
    }
    _ = file // use file here
}
```

#### Example usage:

```go
import "github.com/goaux/iter/scope"

for file, err := range scope.Use2(os.Create("test.gz")) {
    // If err==nil, the file will be closed at the end of the loop body regardless of break.
    if err != nil {
        fmt.Println(err)
        break
    }
    for gz := range scope.Use(gzip.NewWriter(file)) {
        // gz will always be closed at the end of the loop body.
        if _, err := fmt.Fprintln(gz, "hello world"); err != nil {
            fmt.Println(err)
            break
        }
    }
}
```
