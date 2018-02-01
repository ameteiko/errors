package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Wrap :: for a single built-in error :: returns an Appender with just one entry
//
func TestWrapForAnError(t *testing.T) {
	e1 := errors.New("1")

	enchainer := Wrap(e1)

	assert.Len(t, enchainer.GetErrors(), 1)
	assert.Equal(t, e1, enchainer.GetErrors()[0])
	assert.Equal(t, "1", enchainer.Error())
}

//
// Wrap :: for a several built-in errors :: returns an Appender with several entries
//
func TestWrapForSeveralErrors(t *testing.T) {
	e1 := errors.New("1")
	e2 := errors.New("2")
	e3 := errors.New("3")

	enchainer := Wrap(e1, e2, e3)

	assert.Len(t, enchainer.GetErrors(), 3)
	assert.Equal(t, "3 : 2 : 1", enchainer.Error())
}

//
// Wrap :: for a built-in error and an enchainer instance that goes first :: returns an Appender with two entries
//
func TestWrapForABuiltInErrorAndEnchainerErrors(t *testing.T) {
	e1 := NewChain().Append(errors.New("1"))
	e2 := errors.New("2")

	enchainer := Wrap(e1, e2)

	assert.Len(t, enchainer.GetErrors(), 2)
	assert.Equal(t, "2 : 1", enchainer.Error())
}

//
// Wrap :: for a built-in errors and an enchainer instance that goes first :: returns an Appender with several entries
//
func TestWrapForABuiltInErrorsAndEnchainerGoesFirst(t *testing.T) {
	e1 := NewChain().Append(errors.New("1"))
	e2 := errors.New("2")
	e3 := errors.New("3")
	e4 := errors.New("4")

	enchainer := Wrap(e1, e2, e3, e4)

	assert.Len(t, enchainer.GetErrors(), 4)
	assert.Equal(t, "4 : 3 : 2 : 1", enchainer.Error())
}

//
// Wrap :: for a built-in error and an enchainer instance that goes last :: returns an Appender with two entries
//
func TestWrapForABuiltInErrorAndEnchainerThatGoesLast(t *testing.T) {
	e1 := errors.New("1")
	e2 := NewChain().Append(errors.New("2"))

	enchainer := Wrap(e1, e2)

	assert.Len(t, enchainer.GetErrors(), 2)
	assert.Equal(t, "2 : 1", enchainer.Error())
}

//
// Wrap :: for a built-in errors and an enchainer instance that goes last :: returns an Appender with several entries
//
func TestWrapForABuiltInErrorsAndEnchainerThatGoesLast(t *testing.T) {
	e1 := errors.New("1")
	e2 := errors.New("2")
	e3 := errors.New("3")
	e4 := NewChain().Append(errors.New("4"))

	enchainer := Wrap(e1, e2, e3, e4)

	assert.Len(t, enchainer.GetErrors(), 4)
	assert.Equal(t, "4 : 3 : 2 : 1", enchainer.Error())
}

//
// Wrap :: for a built-in errors and an enchainer in between :: returns an Appender with several entries
//
func TestWrapForABuiltInErrorsAndEnchainerInBetween(t *testing.T) {
	e1 := errors.New("1")
	e2 := errors.New("2")
	e3 := errors.New("3")
	e4 := NewChain().Append(errors.New("4"))
	e5 := errors.New("5")

	enchainer := Wrap(e1, e2, e3, e4, e5)

	assert.Len(t, enchainer.GetErrors(), 5)
	assert.Equal(t, "5 : 4 : 3 : 2 : 1", enchainer.Error())
}

//
// Wrap :: for a built-in errors and an several enchainers in between :: returns an Appender with several entries
//
func TestWrapForABuiltInErrorsAndSeveralEnchainersInBetween(t *testing.T) {
	e1 := errors.New("1")
	e2 := errors.New("2")
	e3 := NewChain().Append(errors.New("31")).Append(errors.New("32"))
	e4 := NewChain().Append(errors.New("4"))
	e5 := errors.New("5")

	enchainer := Wrap(e1, e2, e3, e4, e5)

	assert.Len(t, enchainer.GetErrors(), 5)
	assert.Equal(t, "5 : 4 : 32 : 31 : 2 : 1", enchainer.Error())
	assert.Equal(t, e4, enchainer)
}

//
// Wrap :: for a single chainer :: returns a chainer instance
//
func TestWrapForASingleChainer(t *testing.T) {
	chain := NewChain()

	chainedChain := Wrap(chain)

	assert.Equal(t, chain, chainedChain)
	assert.Len(t, chainedChain.GetErrors(), 0)
}

//
// Cause :: for custom error type selected by type reference :: returns custom error
//
func TestCauseForExistingErrorType(t *testing.T) {
	err := NewChain()
	e1 := fmt.Errorf("1")
	e2 := customError{"2"}
	e3 := fmt.Errorf("3")
	e4 := fmt.Errorf("4")

	err.Append(e1).Append(e2).Append(e3).Append(e4)
	e := Cause(err, (*customError)(nil))

	assert.NotEmpty(t, e)
	assert.Equal(t, e2, e)
}

//
// Cause :: for custom error type selected by interface :: returns nil
//
func TestCauseForAnInterface(t *testing.T) {
	err := NewChain()
	e1 := fmt.Errorf("1")
	e2 := customError{"2"}
	e3 := fmt.Errorf("3")
	e4 := fmt.Errorf("4")

	err.Append(e1).Append(e2).Append(e3).Append(e4)
	e := Cause(err, (customMessenger)(nil))

	assert.Nil(t, e)
}

//
// Cause :: for custom error type selected by interface pointer :: returns error
//
func TestCauseForAnInterfaceReference(t *testing.T) {
	err := NewChain()
	e1 := fmt.Errorf("1")
	e2 := customError{"2"}
	e3 := fmt.Errorf("3")
	e4 := fmt.Errorf("4")

	err.Append(e1).Append(e2).Append(e3).Append(e4)
	e := Cause(err, (*customMessenger)(nil))

	assert.NotEmpty(t, e)
	assert.Equal(t, e2, e)
}

//
// Cause :: for several custom error types selected by type reference :: returns the first custom error
//
func TestCauseForSeveralCustomErrors(t *testing.T) {
	err := NewChain()
	e1 := fmt.Errorf("1")
	e2 := customError{"2"}
	e3 := fmt.Errorf("3")
	e4 := customError{"4"}
	e5 := fmt.Errorf("5")

	err.Append(e1).Append(e2).Append(e3).Append(e4).Append(e5)
	e := Cause(err, (*customError)(nil))

	assert.NotEmpty(t, e)
	assert.Equal(t, e4, e)
}

//
// Cause :: for an error that doesn't exist in a chain :: returns nil
//
func TestCauseForAnErrorThatDoesNotExistInChain(t *testing.T) {
	err := errors.New(errMsg)
	chain := Wrap(err, errors.New(infoMsg))

	e := Cause(chain, (*Err)(nil))

	assert.Empty(t, e)
}

//
// Cause :: for a plain error :: returns an error
//
func TestCauseForAPlainError(t *testing.T) {
	err := New(errMsg)

	e := Cause(err, (*Err)(nil))

	assert.Error(t, e)
	assert.Equal(t, e, err)
}

//
// WithMessage :: for an empty message :: returns original error chain
//
func TestWithMessageForAnEmptyMessage(t *testing.T) {
	chain := Wrap(New(errMsg))
	chainErrs := len(chain.GetErrors())

	c := WithMessage(chain, emptyMsg)
	newChain := c.(Chainer)
	newChainErrs := len(newChain .GetErrors())

	assert.Equal(t, chainErrs, newChainErrs)
}

//
// WithMessage :: for an empty message and plaint error (not a chainer) :: original error wrapped into chainer
//
func TestWithMessageForAnEmptyMessageAndPlainError(t *testing.T) {
	err := New(errMsg)

	c := WithMessage(err, emptyMsg)
	newChain := c.(Chainer)
	newChainErrs := len(newChain.GetErrors())

	assert.Equal(t, 1, newChainErrs)
}

//
// WithMessage :: for message and plaint error (not a chainer) :: returns original error wrapped into chainer
//
func TestWithMessageForAMessageAndPlainError(t *testing.T) {
	err := New(errMsg)

	c := WithMessage(err, errMsg)
	newChain := c.(Chainer)
	newChainErrs := len(newChain.GetErrors())

	assert.Equal(t, 2, newChainErrs)
}

//
// WithMessage :: for message and a chainer object :: returns passed chainer with an error.
//
func TestWithMessageForAMessageAndAChainer(t *testing.T) {
	chain := Wrap(New(errMsg))

	c := WithMessage(chain, errMsg)
	newChain := c.(Chainer)
	newChainErrs := len(newChain.GetErrors())

	assert.Equal(t, 2, newChainErrs)
}

//
// customError is a custom error type for tests.
//
type customError struct {
	msg string
}

//
// customMessenger is an interface for a customError.
//
type customMessenger interface {
	GetMessage() string
}

//
// Chain implements error interface.
//
func (e customError) Error () string {

	return e.msg
}

//
// GetMessage implements customMessenger interface
//
func (e customError) GetMessage() string {

	return e.msg
}