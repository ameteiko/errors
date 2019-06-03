package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ErrorMatchesTestSuite is a test suite for all errorMatches() scenarios.
type ErrorMatchesTestSuite struct {
	suite.Suite
	err1 error
}

// TestErrorMatchesTestSuite runs all tests from the ErrorMatchesTestSuite.
func TestErrorMatchesTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorMatchesTestSuite))
}

// SetupSuite initializes the ErrorMatchesTestSuite.
func (s *ErrorMatchesTestSuite) SetupSuite() {
	s.err1 = errors.New("1")
}

// errorMatches() :: for nil values :: returns nil.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForNilParameters() {

	match := errorMatches(nil, nil)

	s.False(match)
}

// errorMatches() :: for two equal errors :: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForTwoSameErrors() {

	match := errorMatches(s.err1, s.err1)

	s.True(match)
}

// errorMatches() :: for target interface pointer not matching source error :: returns false.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForACustomInterfacePointer() {

	match := errorMatches(s.err1, (*customInterface)(nil))

	s.False(match)
}

// errorMatches() :: for an error interface pointer :: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForErrorInterfacePointer() {

	match := errorMatches(s.err1, (*error)(nil))

	s.True(match)
}

// ErrorMatches() :: for a custom error and error interface pointer:: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForACustomErrorAndErrorInterfacePointer() {
	customErr := testError{msg: "error message"}

	match := errorMatches(customErr, (*error)(nil))

	s.True(match)
}

// ErrorMatches() :: for a custom error and custom interface pointer:: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForACustomErrorAndCustomInterfacePointer() {
	customErr := testError{msg: "error message"}

	match := errorMatches(customErr, (*testErrorInterface)(nil))

	s.True(match)
}

// ErrorMatches() :: for a custom error and custom interface :: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForACustomErrorAndCustomInterface() {
	customErr := testError{msg: "error message"}

	match := errorMatches(customErr, (testErrorInterface)(nil))

	s.False(match)
}

// ErrorMatches() :: for a custom error and a pointer to error struct :: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForACustomErrorAndAPointerToErrorStruct() {
	customErr := testError{msg: "error message"}

	match := errorMatches(customErr, (*testError)(nil))

	s.True(match)
}

// ErrorMatches() :: for a pointer to a custom error and a pointer to error struct :: returns true.
func (s *ErrorMatchesTestSuite) TestErrorMatchesForACustomErrorPointerAndAPointerToErrorStruct() {
	customErr := testError{msg: "error message"}

	match := errorMatches(&customErr, (*testError)(nil))

	s.True(match)
}

// customInterface is an interface for testing purposes.
type customInterface interface {
	Send() string
}
