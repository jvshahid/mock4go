package gomock

import (
	// "fmt"
	"reflect"
)

type Function interface{}

var Map = make(map[Function][]*FunctionCall)

type FunctionCall struct {
	args   []interface{}
	values []interface{}
}

func Mock(fun Function, args ...interface{}) *FunctionCall {
	call := &FunctionCall{
		args: args,
	}
	funType := reflect.ValueOf(fun)
	calls := Map[funType]
	calls = append(calls, call)
	Map[funType] = calls
	return call
}

type MethodCall struct {
	FunctionCall
	receiver interface{}
}

func MockWithReceiver(fun Function, receiver interface{}, args ...interface{}) *MethodCall {
	call := &MethodCall{
		FunctionCall: FunctionCall{
			args: args,
		},
		receiver: receiver,
	}
	return call
}

func (m *FunctionCall) Return(values ...interface{}) {
	m.values = values
}

// Returns the (return values, true, nil) if the method/function is mocked
// and the args match the expected values. Otherwise, it returns (nil, true, nil)
// if there was an error this function returns (nil, false, error)
func FunctionCalled(fun Function, args ...interface{}) ([]interface{}, bool, error) {
	funType := reflect.ValueOf(fun)
	calls := Map[funType]
outer:
	for _, call := range calls {
		if len(call.args) > len(args) {
			continue
		}
		for idx, arg := range call.args {
			argType := reflect.TypeOf(args[idx])
			if argType.Kind() == reflect.Ptr {
				if arg != args[idx] {
					continue outer
				}
			}
			if !reflect.DeepEqual(arg, args[idx]) {
				continue outer
			}
		}
		return call.values, true, nil
	}
	// what should we do here
	return nil, false, nil
}

func ResetMocks() {
	Map = make(map[Function][]*FunctionCall)
}
