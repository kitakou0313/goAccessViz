package testpkg

// Simple test functions for ReadGraph testing
func FunctionA() {
	FunctionB()
	FunctionC()
}

func FunctionB() {
	FunctionC()
}

func FunctionC() {
	// Base function with no calls
}

func FunctionD() {
	FunctionA()
}