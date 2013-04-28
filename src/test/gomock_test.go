package test

import (
	"errors"
	. "github.com/jvshahid/gomock"
	. "launchpad.net/gocheck"
	"os"
	"reflect"
	"strings"
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

func (suite *GoMockSuite) TestEnvironment(c *C) {
	c.Assert(os.Getenv("GOMOCK_TEST_ENV"), Equals, "gomock")
}

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
	Mock(func() {
		NoReturnValuesNoReceiver("bar")
	})
	NoReturnValuesNoReceiver("bar")
	c.Assert(noReturnValues, Equals, "foo")
}

func (suite *GoMockSuite) TestMockingFunctionsWithOneReturnValueAndNoReceiver(c *C) {
	Mock(func() {
		When(OneReturnValueNoReceiver()).Return("bar")
	})
	c.Assert(OneReturnValueNoReceiver(), Equals, "bar")
	c.Assert(OneReturnValueNoReceiver2(), Equals, "foo2")
}

func (suite *GoMockSuite) TestMockingFunctionsWithMultipleReturnValuesAndNoReceiver(c *C) {
	expectedErr := errors.New("foobar")
	Mock(func() {
		When(MultipleReturnValuesNoReceiver("bar")).Return("foobar", expectedErr)
	})
	val, err := MultipleReturnValuesNoReceiver("foo")
	c.Assert(err, IsNil)
	c.Assert(val, Equals, "foo")
	val, err = MultipleReturnValuesNoReceiver("bar")
	c.Assert(err, Equals, expectedErr)
	c.Assert(val, Equals, "foobar")
}

type PrefixMatcher struct {
	value string
}

func (m *PrefixMatcher) Matches(other interface{}) bool {
	return strings.HasPrefix(other.(string), m.value)
}

func (suite *GoMockSuite) TestMockingWithMatchers(c *C) {
	expectedErr := errors.New("foobar")
	Mock(func() {
		When(MultipleReturnValuesNoReceiver("")).
			WithMatchers(&PrefixMatcher{value: "ba"}). // ignore the values passed before and use the matcher instead
			Return("foobar", expectedErr)
	})
	val, err := MultipleReturnValuesNoReceiver("foo")
	c.Assert(err, IsNil)
	c.Assert(val, Equals, "foo")
	val, err = MultipleReturnValuesNoReceiver("bar")
	c.Assert(err, Equals, expectedErr)
	c.Assert(val, Equals, "foobar")
	val, err = MultipleReturnValuesNoReceiver("baz")
	c.Assert(err, Equals, expectedErr)
	c.Assert(val, Equals, "foobar")
}

func (suite *GoMockSuite) TestMockingFunctionWithNoReturnValues(c *C) {
	foo := &Foo{Field: ""}
	bar := &Foo{Field: ""}
	Mock(func() {
		bar.NoReturnValues("bar")
	})
	foo.NoReturnValues("foo")
	bar.NoReturnValues("bar")
	c.Assert(foo.Field, Equals, "foo")
	// bar.NoReturnValues was stubbed and shouldn't change the value of Field
	c.Assert(bar.Field, Equals, "")
}

func (suite *GoMockSuite) TestMockingInterface(c *C) {
	mock := &MockTestInterface{}
	Mock(func() {
		When(mock.Value()).Return("foo")
	})
	c.Assert(mock.Value(), Equals, "foo")
}

func (suite *GoMockSuite) TestMockingNoResultInterface(c *C) {
	mock := &MockTestNoResultInterface{}
	Mock(func() {
		mock.Value("foo")
	})
	mock.Value("foo")
}

func (suite *GoMockSuite) TestMockingNoArgNameInterface(c *C) {
	mock := &MockTestNoArgNameInterface{}
	Mock(func() {
		When(mock.Value("foo")).Return("bar")
	})
	c.Assert(mock.Value("foo"), Equals, "bar")
}

func (suite *GoMockSuite) TestMockingEmbeddedInterface(c *C) {
	mock := &MockTestEmbeddedInterface{}
	Mock(func() {
		When(mock.Value()).Return("foo")
		When(mock.AnotherValue()).Return("bar")
	})
	c.Assert(mock.Value(), Equals, "foo")
	c.Assert(mock.AnotherValue(), Equals, "bar")
}
