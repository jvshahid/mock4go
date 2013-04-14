package testc

import (
	. "launchpad.net/gocheck"
	"testing"
)

type Function interface{}

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) {
	TestingT(t)
}

type GoMockSuite struct{}

var _ = Suite(&GoMockSuite{})

func (suite *GoMockSuite) TestSing(c *C) {
	c.Assert(Sin(0.0), Equals, 0.0)
}
