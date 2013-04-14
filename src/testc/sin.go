package testc

// #include <math.h>
// #cgo LDFLAGS: -lm
import "C"

func Sin(f float64) float64 {
	return float64(C.sin(C.double(f)))
}
