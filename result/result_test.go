//go:generate nanotemplate -T string --input=result.go.t
//go:generate nanotemplate -T int --input=result.go.t
package result

import (
	"errors"
	"github.com/nanoservice/monad.go/result/result_int"
	"github.com/nanoservice/monad.go/result/result_string"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringExample(t *testing.T) {
	helloFn := func(name string) result_string.Result {
		return result_string.Success("hello, " + name)
	}

	success := result_string.Success("world").Bind(helloFn)
	assert.Equal(t, result_string.Success("hello, world"), success)

	err := errors.New("The error")
	failure := result_string.Failure(err).Bind(helloFn)
	assert.Equal(t, result_string.Failure(err), failure)
}

func TestIntExample(t *testing.T) {
	addTwo := func(x int) result_int.Result {
		return result_int.Success(2 + x)
	}

	success := result_int.Success(7).Bind(addTwo)
	assert.Equal(t, result_int.Success(9), success)

	err := errors.New("The error")
	failure := result_int.Failure(err).Bind(addTwo)
	assert.Equal(t, result_int.Failure(err), failure)
}

func TestOnErrorFn(t *testing.T) {
	var called bool
	var got error
	var r result_string.Result

	called = false
	r = result_string.Success("yep!").
		OnErrorFn(func(e error) { called = true })
	assert.Equal(t, false, called)
	assert.Equal(t, result_string.Success("yep!"), r)

	called = false
	err := errors.New("The error")
	r = result_string.Failure(err).
		OnErrorFn(func(e error) { called = true; got = e })
	assert.Equal(t, true, called)
	assert.Equal(t, err, got)
	assert.Equal(t, result_string.Failure(err), r)
}

func TestSuccessChain(t *testing.T) {
	r := result_int.Success(15).Chain(
		func(x int) result_int.Result {
			return result_int.Success(x + 2)
		},

		func(x int) result_int.Result {
			return result_int.Success(x * 2)
		},

		func(x int) result_int.Result {
			return result_int.Success(x / 3)
		},
	)

	assert.Equal(t, result_int.Success(11), r)
}

func TestFailedChain(t *testing.T) {
	err := errors.New("The error")

	r := result_string.Success("world").Chain(
		func(name string) result_string.Result {
			return result_string.Success("hello, " + name)
		},

		func(greeting string) result_string.Result {
			return result_string.Failure(err)
		},

		func(_ string) result_string.Result {
			return result_string.Success("bye")
		},
	)

	assert.Equal(t, result_string.Failure(err), r)
}

func TestDeferOnSuccess(t *testing.T) {
	executed := false
	var got int

	result_int.
		Success(35).
		Defer(func(x int) { executed = true; got = x }).
		Err()

	assert.Equal(t, true, executed)
	assert.Equal(t, 35, got)
}

func TestDeferOnFailure(t *testing.T) {
	executed := false
	err := errors.New("The error")

	result_int.
		Failure(err).
		Defer(func(_ int) { executed = true }).
		Err()

	assert.Equal(t, false, executed)
}

func TestDeferIsPreserved(t *testing.T) {
	executed := false
	var got int
	err := errors.New("The error")

	result_int.
		Success(24).
		Defer(func(x int) { executed = true; got = x }).Bind(
		func(x int) result_int.Result {
			return result_int.Failure(err)
		}).
		Err()

	assert.Equal(t, true, executed)
	assert.Equal(t, 24, got)
}

func TestDeferCallToErrIsRequired(t *testing.T) {
	executed := false

	result_int.
		Success(25).
		Defer(func(x int) { executed = true })

	assert.Equal(t, false, executed)
}
