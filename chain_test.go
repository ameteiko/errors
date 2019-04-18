package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Chain.Format :: with %s flag for an empty object :: returns an empty string
//
func TestFormatSFlaggedForAnEmptyErrorObject(t *testing.T) {
	err := NewChain()

	o := fmt.Sprintf("%s", err)

	assert.Empty(t, o)
}

//
// Chain.Format :: with %v flag for an empty object :: returns an empty string
//
func TestFormatVFlaggedForAnEmptyErrorObject(t *testing.T) {
	err := NewChain()

	o := fmt.Sprintf("%v", err)

	assert.Empty(t, o)
}

//
// Chain.Format :: with %s flag for an object with single value :: returns expected error message
//
func TestFormatSFlaggedForASingleErrorObject(t *testing.T) {
	errMsg := "some error"
	e := fmt.Errorf(errMsg)

	err := NewChain()
	err.Append(e)
	o := fmt.Sprintf("%s", err)

	assert.Equal(t, errMsg, o)
}

//
// Chain.Format :: with %s flag for an object with several values :: returns expected error message
//
func TestFormatSFlaggedForAMultivaluedErrorObject(t *testing.T) {
	errMsg1 := "some first error"
	errMsg2 := "some second error"
	e1 := fmt.Errorf(errMsg1)
	e2 := fmt.Errorf(errMsg2)
	expectedError := fmt.Sprintf("%s : %s", errMsg2, errMsg1)

	err := NewChain()
	err.Append(e1)
	err.Append(e2)
	o := fmt.Sprintf("%s", err)

	assert.Equal(t, expectedError, o)
}

//
// Chain.Format :: with %v flag for an object with several values :: returns expected error message
//
func TestFormatVFlaggedForAMultivaluedErrorObject(t *testing.T) {
	errMsg1 := "some first error"
	errMsg2 := "some second error"
	e1 := fmt.Errorf(errMsg1)
	e2 := fmt.Errorf(errMsg2)
	expectedError := fmt.Sprintf("%s : %s", errMsg2, errMsg1)

	err := NewChain()
	err.Append(e1)
	err.Append(e2)
	o := fmt.Sprintf("%v", err)

	assert.Equal(t, expectedError, o)
}

//
// Chain.Append, Chain.prependError :: with different ordering :: returns an expected order
//
func TestAppendPrependError(t *testing.T) {
	err1 := fmt.Errorf("1")
	err2 := fmt.Errorf("2")
	err3 := fmt.Errorf("3")
	err4 := fmt.Errorf("4")
	err5 := fmt.Errorf("5")

	err := NewChain()
	err.Append(err3).Append(err4).Prepend(err2).Prepend(err1).Append(err5)
	o := fmt.Sprintf("%s", err)

	assert.Equal(t, "5 : 4 : 3 : 2 : 1", o)
}

//
// Chain.WithMessage :: with message :: returns an expected error
//
func TestWithMessage(t *testing.T) {
	err1 := fmt.Errorf("1")
	err2 := fmt.Errorf("2")

	err := NewChain().Append(err1).Append(err2).WithMessage("3")
	o := fmt.Sprintf("%s", err)

	assert.Equal(t, "3 : 2 : 1", o)
}
