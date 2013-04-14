package testc

import (
	. "github.com/jvshahid/gomock"
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

func (suite *GoMockSuite) SetUpSuite(c *C) {
	// setup the suite
}

func (suite *GoMockSuite) SetUpTest(c *C) {
	// setup the test
}

func (suite *GoMockSuite) TearDownTest(c *C) {
	ResetMocks()
}

func (suite *GoMockSuite) TearDownSuite(c *C) {
	// tear down the suite
}

func (suite *GoMockSuite) TestSing(c *C) {
	c.Assert(Sin(0.0), Equals, 0.0)
}
