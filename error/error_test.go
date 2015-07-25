package error

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBindHelperAlwaysExecutesProvidedBlock(t *testing.T) {
	executed := false
	Bind(func() error {
		executed = true
		return nil
	})

	assert.Equal(t, true, executed)
}

func TestBindHelperReturnsNonFailedErrorIfNil(t *testing.T) {
	e := Bind(func() error { return nil })
	assert.Equal(t, Return(nil), e)
}

func TestBindHelperReturnsWrappedErrorIfFails(t *testing.T) {
	err := errors.New("Something gone wrong")
	e := Bind(func() error { return err })
	assert.Equal(t, Return(err), e)
}

func TestReturnWrapsNil(t *testing.T) {
	e := Return(nil)
	assert.Equal(t, Error{nil, make([]deferrableFunc, 0)}, e)
}

func TestReturnWrapsError(t *testing.T) {
	err := errors.New("Something else have gone wrong")
	e := Return(err)
	assert.Equal(t, Error{err, make([]deferrableFunc, 0)}, e)
}

func TestBindOnNoErrorExecutesProvidedBlock(t *testing.T) {
	executed := false
	Return(nil).Bind(func() error {
		executed = true
		return nil
	})

	assert.Equal(t, true, executed)
}

func TestBindOnNoErrorReturnsNoErrorIfNotFails(t *testing.T) {
	e := Return(nil).Bind(func() error { return nil })
	assert.Equal(t, Return(nil), e)
}

func TestBindOnNoErrorReturnsWrappedErrorIfFails(t *testing.T) {
	err := errors.New("Unable to parse data")
	e := Return(nil).Bind(func() error { return err })
	assert.Equal(t, Return(err), e)
}

func TestBindOnErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	executed := false
	err := errors.New("Incompatible message version")
	Return(err).Bind(func() error {
		executed = true
		return nil
	})

	assert.Equal(t, false, executed)
}

func TestBindOnErrorReturnsSameError(t *testing.T) {
	err := errors.New("Out of imagination to create new error")
	e := Return(err)
	e2 := e.Bind(func() error { return nil })
	assert.Equal(t, e, e2)
}

func TestDeferOnNoErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	executed := false
	Return(nil).Defer(func() { executed = true })
	assert.Equal(t, false, executed)
}

func TestDeferOnNoErrorReturnsSameValue(t *testing.T) {
	e := Return(nil)
	e2 := e.Defer(func() {})
	assert.Equal(t, e.err, e2.err)
}

func TestDeferOnNoErrorAfterErrExecutesProvidedBlock(t *testing.T) {
	executed := false
	Return(nil).Defer(func() { executed = true }).Err()
	assert.Equal(t, true, executed)
}

func TestDeferOnErrorDoesNotExecuteProvidedBlock(t *testing.T) {
	executed := false
	err := errors.New("Yet Another Error Message (YAEM)")
	Return(err).Defer(func() { executed = true })
	assert.Equal(t, false, executed)
}

func TestDeferOnErrorReturnsSameValue(t *testing.T) {
	err := errors.New("Yet Another YAEM")
	e := Return(err)
	e2 := e.Defer(func() {})
	assert.Equal(t, e, e2)
}

func TestDeferOnErrorAfterErrDoesNotExecuteProvidedBlock(t *testing.T) {
	err := errors.New("Yet Another YAEM")
	executed := false
	Return(err).Defer(func() { executed = true }).Err()
	assert.Equal(t, false, executed)
}

func TestErrOnNoErrorReturnsNil(t *testing.T) {
	assert.Equal(t, nil, Return(nil).Err())
}

func TestErrOnErrorReturnsInnerValue(t *testing.T) {
	err := errors.New("Unable to connect to server")
	assert.Equal(t, err, Return(err).Err())
}

func TestOnErrorOnNoErrorReturnsSpecialError(t *testing.T) {
	e := Return(nil).OnError()
	assert.Equal(t, Return(ErrorWasExpected), e)
}

func TestOnErrorOnErrorReturnsNoError(t *testing.T) {
	err := errors.New("Some Error")
	e := Return(err).OnError()
	assert.Equal(t, Return(nil), e)
}
