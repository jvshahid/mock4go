package test

import (
	"gomock"
	. "launchpad.net/gocheck"
	"reflect"
	"testing"
)

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
	// setup the suite
}

type Function interface{}

func (suite *GoMockSuite) TestBasicAssumptionsAboutFunctions(c *C) {
	var fun Function = NoReturnValuesNoReceiver
	var fun2 Function = NoReturnValuesNoReceiver
	var fun3 Function = NoReturnValuesNoReceiver2
	c.Assert(reflect.ValueOf(fun), Equals, reflect.ValueOf(fun2))
	c.Assert(reflect.ValueOf(fun), Not(Equals), reflect.ValueOf(fun3))
}

func (suite *GoMockSuite) TestBasicMocking(c *C) {
	gomock.Mock(OneReturnValueNoReceiver).Return("bar")
	c.Assert(OneReturnValueNoReceiver(), Equals, "bar")
}

func (suite *GoMockSuite) TestNoMocking(c *C) {
	c.Assert(OneReturnValueNoReceiver(), Equals, "foo")
}
