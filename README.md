## Why another mocking library for GO

When I started writing gomock, there were two other projects
[here](http://godoc.org/code.google.com/p/gomock/gomock) and
[here](https://github.com/jacobsa/gomock) that took a similar approach.
The approach is to manually generate source code for interfaces using a provided
tool.

I didn't like this approach for the following reasons:

1. It requires a manual step of generating source code for each
interfaces that will be used in testing.
2. It only mock interfaces.

I decided to take a different approach after using
[gocov](https://github.com/axw/gocov) and being inspired by their
approach. gomock will create an instrumented copy of the code and
run go test using the new copy. This allows gomock to insert code
to intercept function calls and do some interesting stuff.

## Disclaimer

This library is still a work in progress and has some rough edges. I encourage
you to start using it and contribute back or share your experience in order
to improve it.

## Install

`go get github.com/jvshahid/gomock/gomock`

## Running the tests

`./bin/gomock my_package`

For further help run `./bin/gomock`.

## Usage Example

The examples below use gocheck as the test framework. To stub a
function you call `When(FunctionName([matchers])[.Return(returnValues)]`
where:

1. `matchers` are the argument matchers, currently we only support
   values matchers. A matcher equals the corresponding positional argument
   iff the matcher value wasn't a pointer and the argument equals
   the matcher value using reflect.DeepEquals, or the matcher value
   was a pointer and it is equal to the argument using `==`
2. `returnValues` can be any number of return values (or even omitted)
3. `When` can be omitted if there's no a `Return` clause.
4. Currently gomock doesn't check that the returnValues
   makes sense. That means that your test may compile but panics during
   runtime.
5. Make sure you run ResetMocks() after each test

### Stubbing functions wo a receiver

Given the very simple function below:

```GO
func OneReturnValueNoReceiver() string {
	return "foo"
}
```

We can stub in the test like this:

```GO
package apackage

import (
	. "github.com/jvshahid/gomock"
	. "launchpad.net/gocheck"
	"testing"
)

type GoMockSuite struct{}

var _ = Suite(&GoMockSuite{})

func (suite *GoMockSuite) TearDownTest(c *C) {
	ResetMocks()
}

func (suite *GoMockSuite) TestMockingFunctionsWithOneReturnValueAndNoReceiver(c *C) {
	Mock(func() {
		When(OneReturnValueNoReceiver()).Return("bar")
	})
	c.Assert(OneReturnValueNoReceiver(), Equals, "bar")
	c.Assert(OneReturnValueNoReceiver2(), Equals, "foo2")
}
```

### Stubbing functions with a receiver

```GO
type Foo struct {
	Field string
}

func (f *Foo) NoReturnValues(value string) {
	f.Field = value
}

```

```GO
func (suite *GoMockSuite) TestMockingFunctionWithNoReturnValues(c *C) {
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
```

### Using matchers

gomock allows you to write your own argument matchers.


```GO
func MultipleReturnValuesNoReceiver(value string) (string, error) {
	return value, nil
}
```

```GO
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
```

The reason you have to call the function with a dummy value is that there is no way to reliably compare
function pointers in Go. A preferred way to do this is the following:

```GO
When(FunctionName, NewPrefixMatcher("ba")).Return("something")
```

but this will require the ability to test for "function" equality in gomock which can't be reliably
done. See http://golang.org/doc/go1.html#equality for more information about why
it was decided to remove function equality in Go 1.0.

### Stubbing interfaces

gomock will create a mock implementation for every interface it parses.
If the interface is called `FooInterface` the generated type will be called
`MockFooInterface`, and all the interface's functions will be defined for
that type.

## TODO

* Enhance the documentation of both the code and usage of the library
* Add more matchers, so we can do interesting things like match on a prefix, etc.
* Add a way to pass a new function that decides what to return to the caller
* Ability to exclude certain packages from being instrumented

## Contributing

If you found a bug, want to add a feature:

1. make sure you can run `./bin/test.sh` and it passes on your local machine
2. write a test in `src/test` and make sure it fails
3. make the test pass
4. send me a pr

If you're feeling lazy or don't know how to fix a bug or implement a feature
feel free to open a new issue and I'll make it happen.

## License:

    (The MIT License)

    Copyright (c) 2013 :

    * {John Shahid}[http://github.com/jvshahid]


    Permission is hereby granted, free of charge, to any person obtaining
    a copy of this software and associated documentation files (the
    'Software'), to deal in the Software without restriction, including
    without limitation the rights to use, copy, modify, merge, publish,
    distribute, sublicense, and/or sell copies of the Software, and to
    permit persons to whom the Software is furnished to do so, subject to
    the following conditions:

    The above copyright notice and this permission notice shall be
    included in all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
    EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
    MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
    IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
    CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
    TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
    SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
