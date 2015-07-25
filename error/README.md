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

  }).Err
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

### `(errorMonad.Error) Bind(fn func() error) errorMonad.Error`

Use `(errorMonad.Error) Bind` function to attach new item to a chain. `fn` will
get called if and only if previous chain item haven't returned error.

```go
e.Bind(func() error {
  return doSomethingCausingError()
})
```
