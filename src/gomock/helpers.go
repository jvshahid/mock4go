package gomock

type MockedFunction struct {
	args   []interface{}
	values []interface{}
}

func Mock(fun interface{}, args ...interface{}) *MockedFunction {
	return &MockedFunction{
		args: args,
	}
}

type MockedFunctionWithReceiver struct {
	MockedFunction
	receiver interface{}
}

func MockWithReceiver(fun interface{}, receiver interface{}, args ...interface{}) *MockedFunctionWithReceiver {
	return &MockedFunctionWithReceiver{
		MockedFunction: MockedFunction{
			args: args,
		},
		receiver: receiver,
	}
}

func (m *MockedFunction) Return(values ...interface{}) {
	m.values = values
}

func FunctionCalled(id int, args ...interface{}) []interface{} {
	// what should we do here
	return nil
}
