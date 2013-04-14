package testnomock

// make sure gomock is copied to the temp directory when a package doesn't import
// gomock.

func SomeFunction() {}

// make sure that % in the code don't screw up the printer
func AnotherFunction() int {
	return 5 % 3
}
