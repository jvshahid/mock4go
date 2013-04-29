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

type Mock4goSuite struct{}

var _ = Suite(&Mock4goSuite{})

func (suite *Mock4goSuite) TestSing(c *C) {
	c.Assert(Sin(0.0), Equals, 0.0)
}
