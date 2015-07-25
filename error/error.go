package error

type failableFunc func() error
type deferrableFunc func()

type Error struct {
	err      error
	deferred []deferrableFunc
}

func Return(value error) Error {
	return Error{value, make([]deferrableFunc, 0)}
}

func Bind(fn failableFunc) Error {
	return Return(fn())
}

func (e Error) Bind(fn failableFunc) Error {
	if e.err != nil {
		return e
	}
	return Return(fn())
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

func (e Error) resolveDeferred() {
	for _, fn := range e.deferred {
		fn()
	}
}
