package api

import (
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type Mock4goTestSuite struct {
}

func (suite *Mock4goTestSuite) SetUpTest(c *C) {
}

func (suite *Mock4goTestSuite) TearDownTest(c *C) {
}

var _ = Suite(&Mock4goTestSuite{})

func (s *Mock4goTestSuite) TestInstrumentFile(c *C) {
	c.Fail()
}
