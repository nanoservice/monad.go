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
