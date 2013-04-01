package api

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type GoMockTestSuite struct {
}

func (suite *GoMockTestSuite) SetUpTest(c *C) {
}

func (suite *GoMockTestSuite) TearDownTest(c *C) {
}

var _ = Suite(&GoMockTestSuite{})

func (s *GoMockTestSuite) TestInstrumentFile(c *C) {
	c.Fail()
}
