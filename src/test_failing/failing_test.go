package test_failing

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) {
	TestingT(t)
}

type Mock4goSuite struct{}

var _ = Suite(&Mock4goSuite{})

func (suite *Mock4goSuite) TestFailingTest(c *C) {
	c.Assert(true, Equals, false)
}
