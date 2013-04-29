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

type Mock4goSuite struct{}

var _ = Suite(&Mock4goSuite{})

func (suite *Mock4goSuite) TestAnotherFunction(c *C) {
	c.Assert(AnotherFunction(), Equals, 2)
}
