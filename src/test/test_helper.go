package test

import (
	"fmt"
)

type TestInterface interface {
	Value() string
}

type TestEmbeddedInterface interface {
	TestInterface
	AnotherValue() string
}

type TestInterfaceMethodWithArgs interface {
	Value(firstName, lastName string) string
}

type TestNoResultInterface interface {
	Value(s string)
}

type TestNoArgNameInterface interface {
	Value(string) string
}

type TestInterfaceMethodWithArgsImpl struct{}

func (m *TestInterfaceMethodWithArgsImpl) Value(firstName, lastName string) string {
	return firstName + lastName
}

var _ TestInterfaceMethodWithArgs = new(TestInterfaceMethodWithArgsImpl)

type Foo struct {
	Field string
}

func (f *Foo) NoReturnValues(value string) {
	f.Field = value
}

func (f *Foo) OneReturnValue() string {
	return f.Field
}

func (f *Foo) MultipleReturnValues() (string, error) {
	if f.Field == "" {
		return f.Field, fmt.Errorf("Value is empty")
	}
	return f.Field, nil
}

var noReturnValues string

func NoReturnValuesNoReceiver(value string) {
	noReturnValues = value
}

func NoReturnValuesNoReceiver2(value string) {
	return
}

func OneReturnValueNoReceiver() string {
	return "foo"
}

func OneReturnValueNoReceiver2() string {
	return "foo2"
}

func MultipleReturnValuesNoReceiver(value string) (string, error) {
	return value, nil
}
