package test

import (
	"fmt"
)

type TestInterface interface {
	Value() string
}

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
