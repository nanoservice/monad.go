// result monad
// generated by github.com/nanoservice/monad.go/result
// type: {{T}}
package result_{{t}}

import ({{I}})

type Result struct {
	value *{{T}}
	err   error
}

func Success(value {{T}}) Result {
	return Result{value: &value, err: nil}
}

func Failure(err error) Result {
	return Result{value: nil, err: err}
}

func (r Result) Bind(fn func({{T}}) Result) Result {
	if r.err != nil {
		return r
	}
	return fn(*r.value)
}

func (r Result) OnErrorFn(fn func(error)) Result {
  if r.err != nil {
    fn(r.err)
  }
  return r
}
