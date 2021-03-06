// Code generated by github.com/nanoservice/monad.go/result
// result monad
// type: *os.File
package result_file

import ("os")

type handler           func(*os.File) Result
type errorHandler      func(error)
type deferHandler      func()
type boundDeferHandler func(*os.File)

type Result struct {
        value         **os.File
        err           error
        deferHandlers []deferHandler
}

func NewResult(value *os.File, err error) Result {
        return buildResult(&value, err)
}

func Success(value *os.File) Result {
        return buildResult(&value, nil)
}

func Failure(err error) Result {
        return buildResult(nil, err)
}

func (r Result) Bind(fn handler) Result {
        if r.err != nil {
          return r
        }

        result := fn(*r.value)
        return r.augment(result.value, result.err)
}

func (r Result) Defer(fn boundDeferHandler) Result {
        if r.err != nil {
                return r
        }

        return Result{
                value:         r.value,
                err:           r.err,
                deferHandlers: append(
                        r.deferHandlers,
                        func() { fn(*r.value) },
                ),
        }
}

func (r Result) Err() error {
        for _, fn := range r.deferHandlers {
                fn()
        }
        return r.err
}

func (r Result) Chain(fns... handler) Result {
        for _, fn := range fns {
                r = r.Bind(fn)
        }
        return r
}

func (r Result) OnErrorFn(fn errorHandler) Result {
        if r.err != nil {
                fn(r.err)
        }
        return r
}

func (r Result) augment(value **os.File, err error) (result Result) {
        result = buildResult(value, err)
        result.deferHandlers = r.deferHandlers
        return
}

func buildResult(value **os.File, err error) Result {
        return Result{
                value:         value,
                err:           err,
                deferHandlers: []deferHandler{},
        }
}
