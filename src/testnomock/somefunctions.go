package testnomock

// make sure mock4go is copied to the temp directory when a package doesn't import
// mock4go.

func SomeFunction() {}

// make sure that % in the code don't screw up the printer
func AnotherFunction() int {
	return 5 % 3
}
