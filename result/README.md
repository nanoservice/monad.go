# Result monad

Monad for handling chain of actions that produce either result or an error.
It allows you to make your code much more confident.

Part of [monad.go](https://github.com/nanoservice/monad.go) library.

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
//go:generate nanoinstall -M result -v v1.3.0
//go:generate nanotemplate -T *package.Resource -t resource -I package --input=_result.tt.go
//go:generate nanotemplate -T *package.Status -t status -I package --input=_result.tt.go

// Can be even better if returns `result_resource.Result`
func doStuff() error {
  return openResource().
    Defer(closeResource).
    Bind(func(resource *package.Resource) result_resource.Result {
      err := fetchStatus(resource).
        Bind(codeUsing(resource))
      return result_resource.NewResult(resource, err)
    }).Err()
}

func openResource() result_resource.Result {
  resource, err := package.Connect()
  return result_resource.NewResult(resource, err)
}

func closeResource(resource *package.Resource) {
  resource.Close()
}

func fetchStatus(resource *package.Resource) result_status.Result {
  status, err := resource.Status()
  return result_status.NewResult(status, err)
}

func codeUsing(resource *package.Resource) func(status *package.Status) result_status.Result {
  return func(status *package.Status) result_status.Result {
    // do stuff here
  }
}
```

## Installation

It is installed using `go generate` tool, `github.com/nanoservice/monad.go/nanoinstall` and `github.com/nanoservice/monad.go/nanotemplate`.

Install `nanoinstall` and `nanotemplate` first:

```bash
go get github.com/nanoservice/monad.go/nanoinstall
go get github.com/nanoservice/monad.go/nanotemplate
```

Then install recent version of `result` monad by adding this to the top of any of your source files:

```go
//go:generate nanoinstall -M result -v v1.3.0
```

To generate your first result monad instance (with concrete type), add this:

```go
//go:generate nanotemplate -T int --input=_result.tt.go
```

Then run `go generate` and you will get these files:

```bash
./
  _result.tt.go
  result_int/
             result_int.t.go
```

## Usage

### `NewResult(value T, err error) Result<T>`

`NewResult(value, err)` constructs new `Result` monad instance given you already have value-error pair.

```go
func OpenResource(config Resource.Config) result_resource.Result {
  resource, err := Resource.Open(config)
  return result_resource.NewResult(resource, err)
}
```

### `Success(value T) Result<T>`

`Success(value)` constructs successful `Result` monad instance given a value.

```go
result_string.Success("a string value")
// => Result <string> {"a string value"}
```

### `Failure(err error) Result<T>`

`Failure(err)` constructs failed `Result` monad instance given an error value.

```go
result_int.Failure(errors.New("Incompatible numbers"))
// => Result <int> {err: Error{"Incompatible numbers"}}
```

### `(Result<T>) Bind(fn func(T) Result<T>) Result<T>`

`Result.Bind(fn)` will call `fn` with value it holds if it is in `Success` state; it will return whatever `fn` returned in that case.

In case monad is in `Failure` state, `Result.Bind(fn)` will not call `fn` and return itself immediately.

```go
result_string.
  Success("world").
  Bind(func(name string) result_string.Result {
    return result_string.Success("hello, " + name)
  })
// => Result <string> {"hello, world"}

result_string.
  Failure(errors.New("Unable to fetch user name")).
  Bind(func(name string) result_string.Result {
    return result_string.Success("hello, " + name)
  })
// => Result <string> {err: Error{"Unable to fetch user name"}}
```

### `(Result<T>) Defer(fn func(T)) Result<T>`

`Result.Defer(fn)` will schedule deferred call to `fn` if it is in `Success` state; it will return itself immediately.

In case monad is in `Failure` state, `Result.Defer(fn)` will not schedule call to `fn` and return itself immediately.

Scheduled functions are executed in the order they are scheduled upon call to `Result.Err()`, regardless of state of the monad instance.

```go
openResource().Defer(func(resource *resource.Resource) {
  resource.Close()
}).Bind(...)
```

### `(Result<T>) Chain(fns... func(T) Result<T>) Result<T>`

`Result.Chain(fns)` is a syntactic sugar for a chain of subsequent `.Bind(fn)` calls.

```go
openResource().Chain(
  fetchMetaConfig,
  connectToBrokers,
  setupListeners,
  setupPublishers,
)
```

### `(Result<T>) OnErrorFn(fn func(error)) Result<T>`

`Result.OnErrorFn(fn)` calls `fn` with error it contains if it is in `Failure` state; returns itself afterwards.

If it is in `Success` state, then `fn` is not called and it returns itself immediately.

```go
openResource().Chain(
  ...
).OnErrorFn(reportBrokenResource)
```

### `(Result<T>) Err() error`

`Result.Err()` fetches error value from the monad instance. Returns `nil` for monad instance in `Success` state.

Intended to be used at the end of monad call chains and closer to the top of the application.

Executes all scheduled functions.

```go
result_int.Success(4).Err()
// => nil

result_int.
  Failure(errors.New("Unable to read configuration file")).
  Err()
// => Error{"Unable to read configuration file"}
```

---

[List of Monads](https://github.com/nanoservice/monad.go#monads)
