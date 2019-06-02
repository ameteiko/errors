package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ErrorsTestSuite is a test suite for exposed functions.
type ErrorsTestSuite struct {
	suite.Suite
	err1      error
	err2      error
	err3      error
	err31     error
	err32     error
	err4      error
	err5      error
	customErr testError
}

// TestErrorsTestSuite runs all ErrorsTestSuite tests.
func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

// SetupSuite initializes the ErrorsTestSuite.
func (s *ErrorsTestSuite) SetupSuite() {
	s.err1 = fmt.Errorf("1")
	s.err2 = fmt.Errorf("2")
	s.err3 = fmt.Errorf("3")
	s.err31 = fmt.Errorf("31")
	s.err32 = fmt.Errorf("32")
	s.err4 = fmt.Errorf("4")
	s.err5 = fmt.Errorf("5")
	s.customErr = testError{"custom error"}
}

// Wrap :: for a single error :: returns a chain with just one entry.
func (s *ErrorsTestSuite) TestWrapForAnError() {

	err := Wrap(s.err1)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 1)
	s.Equal(s.err1, chain.getErrors()[0])
	s.Equal("1", chain.Error())
}

// Wrap :: for a several errors :: returns a chain with several entries.
func (s *ErrorsTestSuite) TestWrapForSeveralErrors() {

	err := Wrap(s.err1, s.err2, s.err3)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 3)
	s.Equal("3 : 2 : 1", chain.Error())
}

// Wrap :: for a chain instance that goes first and an error :: returns a chain with two entries.
func (s *ErrorsTestSuite) TestWrapForAChainInstanceAndAnError() {
	c1 := newChain()
	c1.append(s.err1)

	err := Wrap(c1, s.err2)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 2)
	s.Equal("2 : 1", chain.Error())
	s.Equal(c1, chain)
	s.Equal(s.err2, chain.getErrors()[0])
	s.Equal(s.err1, chain.getErrors()[1])
}

// Wrap :: for a a enchain instance that goes first and three errors :: returns a chain with several entries.
func (s *ErrorsTestSuite) TestWrapForAEnchainInstanceAndThreeErrors() {
	c1 := newChain()
	c1.append(s.err1)

	err := Wrap(c1, s.err2, s.err3, s.err4)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 4)
	s.Equal("4 : 3 : 2 : 1", chain.Error())
}

// Wrap :: for an error and a chain instance that goes last :: returns a chain with two entries.
func (s *ErrorsTestSuite) TestWrapForAnErrorAndChainInstanceThatGoesLast() {
	c2 := newChain()
	c2.append(s.err2)

	err := Wrap(s.err1, c2)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 2)
	s.Equal("2 : 1", chain.Error())
	s.Equal(c2, chain)
	s.Equal(s.err2, chain.getErrors()[0])
	s.Equal(s.err1, chain.getErrors()[1])
}

// Wrap :: for three errors and a chain instance that goes last :: returns a chain with several entries.
func (s *ErrorsTestSuite) TestWrapForAThreeErrorsAndChainInstanceThatGoesLast() {
	c4 := newChain()
	c4.append(s.err4)

	err := Wrap(s.err1, s.err2, s.err3, c4)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 4)
	s.Equal("4 : 3 : 2 : 1", chain.Error())
}

// Wrap :: for a errors and a chain instance in between :: returns a chain with several entries.
func (s *ErrorsTestSuite) TestWrapForErrorsAndChainInstanceInBetween() {
	c4 := newChain()
	c4.append(s.err4)

	err := Wrap(s.err1, s.err2, s.err3, c4, s.err5)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 5)
	s.Equal("5 : 4 : 3 : 2 : 1", chain.Error())
}

// Wrap :: for errors and an several chain instances in between :: returns a chain with several entries.
func (s *ErrorsTestSuite) TestWrapForABuiltInErrorsAndSeveralEnchainersInBetween() {
	c3 := newChain()
	c3.append(s.err31)
	c3.append(s.err32)
	c4 := newChain()
	c4.append(s.err4)

	err := Wrap(s.err1, s.err2, c3, c4, s.err5)
	chain := err.(*chain)

	s.Len(chain.getErrors(), 6)
	s.Equal("5 : 4 : 32 : 31 : 2 : 1", chain.Error())
	s.Equal(c4, chain)
}

// Wrap :: for a single chain instance :: returns a chain instance.
func (s *ErrorsTestSuite) TestWrapForASingleChainer() {
	c := newChain()

	err := Wrap(c)
	chain := err.(*chain)

	s.Equal(c, chain)
	s.Len(chain.getErrors(), 0)
}

// Cause :: for custom error type selected by type reference :: returns custom error.
func (s *ErrorsTestSuite) TestCauseForExistingErrorType() {
	err2 := testError{"2"}
	c := newChain()
	c.append(s.err1)
	c.append(err2)
	c.append(s.err3)
	c.append(s.err4)

	err := Cause(c, (*testError)(nil))

	s.NotEmpty(err)
	s.Equal(err2, err)
}

// Cause :: for custom error type selected by interface :: returns nil.
func (s *ErrorsTestSuite) TestCauseForAnInterface() {
	c := newChain()
	c.append(s.err1)
	c.append(testError{"2"})
	c.append(s.err3)
	c.append(s.err4)

	err := Cause(c, (testErrorInterface)(nil))

	s.Nil(err)
}

// Cause :: for custom error type selected by interface pointer :: returns error.
func (s *ErrorsTestSuite) TestCauseForAnInterfaceReference() {
	err2 := testError{"2"}
	c := newChain()
	c.append(s.err1)
	c.append(err2)
	c.append(s.err3)
	c.append(s.err4)

	err := Cause(c, (*testErrorInterface)(nil))

	s.NotEmpty(err)
	s.Equal(err2, err)
}

// Cause :: for several custom error types selected by type reference :: returns the last appended custom error.
func (s *ErrorsTestSuite) TestCauseForSeveralCustomErrors() {
	e2 := testError{"2"}
	e4 := testError{"4"}
	c := newChain()
	c.append(s.err1)
	c.append(e2)
	c.append(s.err3)
	c.append(e4)
	c.append(s.err5)

	err := Cause(c, (*testError)(nil))

	s.NotEmpty(err)
	s.Equal(e4, err)
}

// Cause :: for an error that doesn't exist in a chain :: returns nil.
func (s *ErrorsTestSuite) TestCauseForAnErrorThatDoesNotExistInChain() {
	chain := Wrap(s.err1, s.err2)

	e := Cause(chain, (*testError)(nil))

	s.Empty(e)
}

// Cause :: for a plain error :: returns an error.
func (s *ErrorsTestSuite) TestCauseForAPlainError() {

	err := Cause(s.customErr, (*testError)(nil))

	s.Error(err)
	s.Equal(s.customErr, err)
}

// WithMessage :: for an empty message :: returns original error chain.
func (s *ErrorsTestSuite) TestWithMessageForAnEmptyMessage() {
	err := Wrap(s.customErr)
	c := err.(*chain)
	cErrs := len(c.getErrors())

	err = WithMessage(c, "")
	newC := err.(*chain)
	newCErrs := len(newC.getErrors())

	s.Equal(cErrs, newCErrs)
}

// WithMessage :: for an empty message and an error (not a chainer) :: returns a chain with the error.
func (s *ErrorsTestSuite) TestWithMessageForAnErrorAndAnEmptyMessage() {

	err := WithMessage(s.customErr, "")
	c := err.(*chain)
	cErrs := len(c.getErrors())

	s.Equal(1, cErrs)
}

// WithMessage :: for message and an error (not a chainer) :: returns a chain with two errors.
func (s *ErrorsTestSuite) TestWithMessageAnErrorAndAMessage() {

	err := WithMessage(s.customErr, "error message")
	c := err.(*chain)
	cErrs := len(c.getErrors())

	s.Equal(2, cErrs)
}

// WithMessage :: for message and a chain object :: returns passed chain instance with a new error for the message.
func (s *ErrorsTestSuite) TestWithMessageForAMessageAndAChainInstance() {
	err := Wrap(s.customErr)

	err = WithMessage(err, "error message")
	c := err.(*chain)
	cErrs := len(c.getErrors())

	s.Equal(2, cErrs)
}

// testErrorInterface is an interface for a testError.
type testErrorInterface interface {
	Message() string
}

// testError is a custom test error type.
type testError struct{ msg string }

// Error implements error interface.
func (e testError) Error() string { return e.msg }

// Message implements testErrorInterface.
func (e testError) Message() string { return e.msg }
