package test

import (
	"errors"
	gomock "github.com/jvshahid/gomock/api"
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
	// setup the test
}

func (suite *GoMockSuite) TearDownTest(c *C) {
	gomock.ResetMocks()
}

func (suite *GoMockSuite) TearDownSuite(c *C) {
	// tear down the suite
}

type Function interface{}

func (suite *GoMockSuite) TestNoMocking(c *C) {
	c.Assert(OneReturnValueNoReceiver(), Equals, "foo")
}

func (suite *GoMockSuite) TestBasicAssumptionsAboutFunctions(c *C) {
	var fun Function = NoReturnValuesNoReceiver
	var fun2 Function = NoReturnValuesNoReceiver
	var fun3 Function = NoReturnValuesNoReceiver2
	c.Assert(reflect.ValueOf(fun), Equals, reflect.ValueOf(fun2))
	c.Assert(reflect.ValueOf(fun), Not(Equals), reflect.ValueOf(fun3))
}

func (suite *GoMockSuite) TestMockingFunctionWithNoReturnValuesAndNoReceiver(c *C) {
	NoReturnValuesNoReceiver("foo")
	c.Assert(noReturnValues, Equals, "foo")
	gomock.Mock(NoReturnValuesNoReceiver, "bar")
	NoReturnValuesNoReceiver("bar")
	c.Assert(noReturnValues, Equals, "foo")
}

func (suite *GoMockSuite) TestMockingFunctionsWithOneReturnValueAndNoReceiver(c *C) {
	gomock.Mock(OneReturnValueNoReceiver).Return("bar")
	c.Assert(OneReturnValueNoReceiver(), Equals, "bar")
	c.Assert(OneReturnValueNoReceiver2(), Equals, "foo2")
}

func (suite *GoMockSuite) TestMockingFunctionsWithMultipleReturnValuesAndNoReceiver(c *C) {
	expectedErr := errors.New("foobar")
	gomock.Mock(MultipleReturnValuesNoReceiver, "bar").Return("foobar", expectedErr)
	val, err := MultipleReturnValuesNoReceiver("foo")
	c.Assert(err, IsNil)
	c.Assert(val, Equals, "foo")
	val, err = MultipleReturnValuesNoReceiver("bar")
	c.Assert(err, Equals, expectedErr)
	c.Assert(val, Equals, "foobar")
}

func (suite *GoMockSuite) TestMockingFunctionWithNoReturnValues(c *C) {
	foo := &Foo{Field: ""}
	bar := &Foo{Field: ""}
	gomock.Mock((*Foo).NoReturnValues, bar)
	foo.NoReturnValues("foo")
	bar.NoReturnValues("bar")
	c.Assert(foo.Field, Equals, "foo")
	c.Assert(bar.Field, Equals, "")
}
