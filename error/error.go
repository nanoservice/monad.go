package error

import "errors"

type failableFunc func() error
type deferrableFunc func()
type handlerFunc func(error)

type Error struct {
	err      error
	deferred []deferrableFunc
}

var ErrorWasExpected = errors.New("Error was expected")

func Return(value error) Error {
	return Error{value, make([]deferrableFunc, 0)}
}

func Bind(fn failableFunc) Error {
	return Return(fn())
}

func Chain(fns ...failableFunc) Error {
	return Return(nil).Chain(fns...)
}

func (e Error) Bind(fn failableFunc) Error {
	if e.err != nil {
		return e
	}
	return e.modify(fn())
}

func (e Error) Chain(fns ...failableFunc) (result Error) {
	result = e
	for _, fn := range fns {
		result = result.Bind(fn)
	}
	return
}

func (e Error) Defer(fn deferrableFunc) Error {
	if e.err != nil {
		return e
	}
	return Error{e.err, append(e.deferred, fn)}
}

func (e Error) Err() error {
	e.resolveDeferred()
	return e.err
}

func (e Error) OnError() Error {
	if e.err == nil {
		return e.modify(ErrorWasExpected)
	}
	return e.modify(nil)
}

func (e Error) OnErrorFn(fn handlerFunc) Error {
	return e.OnError().Bind(func() error {
		fn(e.err)
		return nil
	})
}

func (e Error) resolveDeferred() {
	for _, fn := range e.deferred {
		fn()
	}
}

func (e Error) modify(err error) Error {
	return Error{err, e.deferred}
}
