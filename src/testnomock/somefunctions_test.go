package testnomock

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

func (suite *GoMockSuite) TestAnotherFunction(c *C) {
	c.Assert(AnotherFunction(), Equals, 2)
}
