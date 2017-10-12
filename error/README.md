# Error monad

Monad for handling chain of actions that can produce errors. In case current
chain item returns error, nullifies the rest of the chain and returns occured
error.

Part of [`monad.go`](https://github.com/nanoservice/monad.go) library.

## Example

Given idiomatic Golang error handling example:

```go
func doStuff() error {
  resource, err := connectResource()
  if err != nil {
    return err
  }

  defer resource.Close()

  status, err := resource.Status()
  if err != nil {
    return err
  }

  return codeUsing(resource, status)
}
```

Can be rewritten as:

```go
func doStuff() error {
  var (
    resource Resource
    status Status
  )

  return errorMonad.Bind(func() (err error) {
    resource, err = connectResource()
    return

  }).Defer(func() {
    resource.Close()

  }).Bind(func() (err error) {
    status, err = resource.Status()
    return

  }).Bind(func() error {
    return codeUsing(resource, status)

  }).Err()
}
```

## Usage

```go
import errorMonad "github.com/nanoservice/monad.go/error"
```

### `errorMonad.Bind(fn func() error) errorMonad.Error`

Use `errorMonad.Bind` function to start the chain. `errorMonad.Bind` will
unconditionally call provided `fn` function and wrap it into monadic context.

```go
errorMonad.Bind(func() error) {
  return doSomethingCausingError()
})
```

### `errorMonad.Return(err error) errorMonad.Error`

Use `errorMonad.Return` function to wrap `err` value in monadic context.

```go
errorMonad.Return(errors.New("Unable to load config"))
```

### `(errorMonad.Error) Chain(fn (func() error)...) errorMonad.Error`

Use `(errorMonad.Error) Chain` function if you find yourself chaining too much
`.Bind(fn)` calls in a row. It has analogous helper function `errorMonad.Chain(fn)`

Given this code:

```go
e.Bind(func() error {
  return doSomething(withThis)

}).Bind(func() error {
  return doSomething(withThat)

}).Bind(func() (err error) {
  result, err = andCalculateTheResult()
  return
})
```

Can be rewritten as:

```go
e.Chain(
  func() error { return doSomething(withThis) },
  func() error { return doSomething(withThat) },
  func() (err error) {
    result, err = andCalculateTheResult()
    return
  },
)
```

### `(errorMonad.Error) Defer(fn func()) errorMonad.Error`

Use `(errorMonad.Error) Defer` function to attach deferred item to a chain.
This item will be enqued for execution if and only if previous chain item
haven't returned error. This is useful for freeing resources after successfully
acquiring them.

Deferred chain items gets executed on call to `(errorMonad.Err) Err()`.

```go
e.Defer(func() {
  resource.Close()
})
```

### `(errorMonad.Error) Err() error`

Use `(errorMonad.Error) Err` function to fetch the error, that failed the
chain. `Err` returns `nil` if chain was successful.

When `Err` is called, all deferred chain items get executed.

```go
e.Err()
```

### `(errorMonad.Error) OnError() errorMonad.Error`

Use `(errorMonad.Error) OnError` to continue chain if and only if there was an
error.

```go
Bind(func() error {
  return ICanFail()
}).OnError().Bind(func() error {
  return ScheduleRetry()
})
```

This function has callback variant `(errorMonad.Error) OnErrorFn(fn func(err error)) errorMonad.Error`:

```go
e.OnErrorFn(func(err error) {
  log.Printf("Some strange error occurred: %v\n", err)
})
```

---

[List of Monads](https://github.com/nanoservice/monad.go#monads)
