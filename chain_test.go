package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ChainTestSuite is a test suite error chain.
type ChainTestSuite struct {
	suite.Suite
	err1         error
	err2         error
	err3         error
	err4         error
	err5         error
	stackMessage string
}

// TestChainTestSuite runs tests from the ChainTestSuite.
func TestChainTestSuite(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}

// SetupSuite performs the suite initialization.
func (s *ChainTestSuite) SetupSuite() {
	s.err1 = fmt.Errorf("1")
	s.err2 = fmt.Errorf("2")
	s.err3 = fmt.Errorf("3")
	s.err4 = fmt.Errorf("4")
	s.err5 = fmt.Errorf("5")
	s.stackMessage = "error stack trace"
}

// append(), prepend() :: with different ordering :: returns an expected order.
func (s *ChainTestSuite) TestAppendPrependError() {
	c := newChain()
	c.append(s.err3)
	c.append(s.err4)
	c.prepend(s.err2)
	c.prepend(s.err1)
	c.append(s.err5)

	o := fmt.Sprintf("%s", c)

	s.Equal("5 : 4 : 3 : 2 : 1", o)
}

// Error() :: for errors with empty messages :: returns empty placeholders for them.
func (s *ChainTestSuite) TestErrorForErrorsWithEmptyMessagesReturnsAPlaceholders() {
	c := newChain()
	c.append(errors.New(""))
	c.append(errors.New(""))

	o := c.Error()

	s.Equal(" : ", o)
	s.Len(c.getErrors(), 2)
}

// Error() :: from some errors with empty messages :: returns empty placeholders for empty messages.
func (s *ChainTestSuite) TestErrorForErrorsWithEmptyMessagesReturnsPlaceholders() {
	c := newChain()
	c.append(errors.New(""))
	c.append(s.err1)
	c.append(errors.New(""))
	c.append(s.err2)

	o := c.Error()

	s.Equal("2 :  : 1 : ", o)
	s.Len(c.getErrors(), 4)
}

// chain.Format :: with %s flag for an empty chain :: returns an empty string.
func (s *ChainTestSuite) TestFormatSFlaggedForAnEmptyChain() {
	c := newChain()

	o := fmt.Sprintf("%s", c)

	s.Empty(o)
}

// chain.Format :: with %v flag for an empty chain :: returns an empty string.
func (s *ChainTestSuite) TestFormatVFlaggedForAnEmptyChain() {
	c := newChain()

	o := fmt.Sprintf("%v", c)

	s.Empty(o)
}

// chain.Format :: with %s flag for a chain with a single value :: returns expected error message.
func (s *ChainTestSuite) TestFormatSFlaggedForASingleErrorObject() {
	c := newChain()
	c.append(s.err1)

	o := fmt.Sprintf("%s", c)

	s.Equal("1", o)
}

// chain.Format :: with %s flag for a chain with several values :: returns expected error message.
func (s *ChainTestSuite) TestFormatSFlaggedForAMultivaluedErrorObject() {
	c := newChain()
	c.append(s.err1)
	c.append(s.err2)

	o := fmt.Sprintf("%s", c)

	s.Equal("2 : 1", o)
}

// chain.Format :: with %v flag for an object with several values :: returns expected error message.
func (s *ChainTestSuite) TestFormatVFlaggedForAMultivaluedErrorObject() {
	c := newChain()
	c.append(s.err1)
	c.append(s.err2)

	o := fmt.Sprintf("%v", c)

	s.Equal("2 : 1", o)
}

// Format() :: for a %+v flag :: returns an error stack trace.
func (s *ChainTestSuite) TestFormatForVFlag() {
	c := newChain()
	c.stack = &formatterStub{s.stackMessage}
	c.append(s.err1)
	c.append(s.err2)

	o := fmt.Sprintf("%+v", c)

	s.Equal("2 : 1\n"+s.stackMessage, o)
}

// formatterStub is a stub for the stack instance to test error formatting.
type formatterStub struct{ msg string }

// Format returns the stub message.
func (fs *formatterStub) Format(state fmt.State, _ rune) { state.Write([]byte(fs.msg)) }
