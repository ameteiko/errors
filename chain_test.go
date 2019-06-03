package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ChainTestSuite is a test suite for all errChain test errors.
type ChainTestSuite struct {
	suite.Suite
	err1 error
	err2 error
	err3 error
	err4 error
	err5 error
}

// TestChainTestSuite runs all the tests from the ChainTestSuite.
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
}

// errChain.Format :: with %s flag for an empty errChain :: returns an empty string.
func (s *ChainTestSuite) TestFormatSFlaggedForAnEmptyChain() {
	c := newChain()

	o := fmt.Sprintf("%s", c)

	s.Empty(o)
}

// errChain.Format :: with %v flag for an empty object :: returns an empty string.
func (s *ChainTestSuite) TestFormatVFlaggedForAnEmptyChain() {
	c := newChain()

	o := fmt.Sprintf("%v", c)

	s.Empty(o)
}

// errChain.Format :: with %s flag for a errChain with a single value :: returns expected error message.
func (s *ChainTestSuite) TestFormatSFlaggedForASingleErrorObject() {
	c := newChain()
	c.append(s.err1)

	o := fmt.Sprintf("%s", c)

	s.Equal("1", o)
}

// errChain.Format :: with %s flag for a errChain with several values :: returns expected error message.
func (s *ChainTestSuite) TestFormatSFlaggedForAMultivaluedErrorObject() {
	c := newChain()
	c.append(s.err1)
	c.append(s.err2)

	o := fmt.Sprintf("%s", c)

	s.Equal("2 : 1", o)
}

// errChain.Format :: with %v flag for an object with several values :: returns expected error message.
func (s *ChainTestSuite) TestFormatVFlaggedForAMultivaluedErrorObject() {
	c := newChain()
	c.append(s.err1)
	c.append(s.err2)

	o := fmt.Sprintf("%v", c)

	s.Equal("2 : 1", o)
}

// errChain.append, errChain.prepend :: with different ordering :: returns an expected order.
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
