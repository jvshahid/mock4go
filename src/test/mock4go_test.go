package test

import (
	"errors"
	. "github.com/jvshahid/mock4go"
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

type Mock4goSuite struct{}

var _ = Suite(&Mock4goSuite{})

func (suite *Mock4goSuite) SetUpSuite(c *C) {
	// setup the suite
}

func (suite *Mock4goSuite) SetUpTest(c *C) {
	// setup the test
}

func (suite *Mock4goSuite) TearDownTest(c *C) {
	ResetMocks()
}

func (suite *Mock4goSuite) TearDownSuite(c *C) {
	// tear down the suite
}

func (suite *Mock4goSuite) TestEnvironment(c *C) {
	c.Assert(os.Getenv("MOCK4GO_TEST_ENV"), Equals, "mock4go")
}

func (suite *Mock4goSuite) TestNoMocking(c *C) {
	c.Assert(OneReturnValueNoReceiver(), Equals, "foo")
}

func (suite *Mock4goSuite) TestBasicAssumptionsAboutFunctions(c *C) {
	var fun Function = NoReturnValuesNoReceiver
	var fun2 Function = NoReturnValuesNoReceiver
	var fun3 Function = NoReturnValuesNoReceiver2
	c.Assert(reflect.ValueOf(fun), Equals, reflect.ValueOf(fun2))
	c.Assert(reflect.ValueOf(fun), Not(Equals), reflect.ValueOf(fun3))
}

func (suite *Mock4goSuite) TestMockingFunctionWithNoReturnValuesAndNoReceiver(c *C) {
	NoReturnValuesNoReceiver("foo")
	c.Assert(noReturnValues, Equals, "foo")
	Mock(func() {
		NoReturnValuesNoReceiver("bar")
	})
	NoReturnValuesNoReceiver("bar")
	c.Assert(noReturnValues, Equals, "foo")
}

func (suite *Mock4goSuite) TestMockingFunctionsWithOneReturnValueAndNoReceiver(c *C) {
	Mock(func() {
		When(OneReturnValueNoReceiver()).Return("bar")
	})
	c.Assert(OneReturnValueNoReceiver(), Equals, "bar")
	c.Assert(OneReturnValueNoReceiver2(), Equals, "foo2")
}

func (suite *Mock4goSuite) TestMockingFunctionsWithMultipleReturnValuesAndNoReceiver(c *C) {
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

func (suite *Mock4goSuite) TestMockingWithMatchers(c *C) {
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

func (suite *Mock4goSuite) TestMockingFunctionWithNoReturnValues(c *C) {
	foo := &Foo{Field: ""}
	bar := &Foo{Field: ""}
	Mock(func() {
		bar.NoReturnValues("bar") // cause the function to be a no operation
	})
	foo.NoReturnValues("foo")
	bar.NoReturnValues("bar")
	c.Assert(foo.Field, Equals, "foo")
	// bar.NoReturnValues was stubbed and shouldn't change the value of Field
	c.Assert(bar.Field, Equals, "")
}

func (suite *Mock4goSuite) TestMockingInterface(c *C) {
	mock := &MockTestInterface{}
	Mock(func() {
		When(mock.Value()).Return("foo")
	})
	c.Assert(mock.Value(), Equals, "foo")
}

func (suite *Mock4goSuite) TestMockingNoResultInterface(c *C) {
	mock := &MockTestNoResultInterface{}
	Mock(func() {
		mock.Value("foo")
	})
	mock.Value("foo")
}

func (suite *Mock4goSuite) TestMockingNoArgNameInterface(c *C) {
	mock := &MockTestNoArgNameInterface{}
	Mock(func() {
		When(mock.Value("foo")).Return("bar")
	})
	c.Assert(mock.Value("foo"), Equals, "bar")
}

func (suite *Mock4goSuite) TestMockingEmbeddedInterface(c *C) {
	mock := &MockTestEmbeddedInterface{}
	Mock(func() {
		When(mock.Value()).Return("foo")
		When(mock.AnotherValue()).Return("bar")
	})
	c.Assert(mock.Value(), Equals, "foo")
	c.Assert(mock.AnotherValue(), Equals, "bar")
}
