package api

import (
	"reflect"
)

type function interface{}

var Map = make(map[function][]*functionCall)

type functionCall struct {
	args   []interface{}
	values []interface{}
}

func getFunType(fun function) interface{} {
	return reflect.ValueOf(fun)
}

var mocking = false
var lastFunctionCall *functionCall

func Mock(fun func()) {
	mocking = true
	fun()
	mocking = false
}

func When(args ...interface{}) *functionCall {
	return lastFunctionCall
}

func (m *functionCall) Return(values ...interface{}) *functionCall {
	m.values = values
	return m
}

func addFunctionCall(funType interface{}, call *functionCall) {
	Map[funType] = append(Map[funType], call)
}

func ZeroValues(fun function) []interface{} {
	funType := reflect.TypeOf(fun)
	values := make([]interface{}, 0)
	for i := 0; i < funType.NumOut(); i++ {
		values = append(values, reflect.Zero(funType.Out(i)).Interface())
	}
	return values
}

// Returns the (return values, true, nil) if the method/function is mocked
// and the args match the expected values. Otherwise, it returns (nil, true, nil)
// if there was an error this function returns (nil, false, error)
func FunctionCalled(fun function, args ...interface{}) ([]interface{}, bool, error) {
	funType := getFunType(fun)
	if mocking {
		lastFunctionCall = &functionCall{
			args: args,
		}
		addFunctionCall(funType, lastFunctionCall)
		return ZeroValues(fun), true, nil
	}
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
	Map = make(map[function][]*functionCall)
}
