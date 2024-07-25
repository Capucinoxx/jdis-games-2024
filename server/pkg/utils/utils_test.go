package utils

import (
	"testing"
)

func TestNilIf(t *testing.T) {
	// Test with int type
	intValue := 5
	if result := NilIf(&intValue, true); result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
	if result := NilIf(&intValue, false); result != &intValue {
		t.Errorf("Expected %v, got %v", &intValue, result)
	}

	// Test with string type
	strValue := "hello"
	if result := NilIf(&strValue, true); result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
	if result := NilIf(&strValue, false); result != &strValue {
		t.Errorf("Expected %v, got %v", &strValue, result)
	}

	// Test with struct type
	type myStruct struct {
		Field int
	}
	structValue := myStruct{Field: 10}
	if result := NilIf(&structValue, true); result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
	if result := NilIf(&structValue, false); result != &structValue {
		t.Errorf("Expected %v, got %v", &structValue, result)
	}

	// Test with nil pointer
	var nilPointer *int
	if result := NilIf(nilPointer, true); result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
	if result := NilIf(nilPointer, false); result != nilPointer {
		t.Errorf("Expected %v, got %v", nilPointer, result)
	}
}
