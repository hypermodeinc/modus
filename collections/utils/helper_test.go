package utils

import (
	"reflect"
	"testing"
)

func TestConvertToFloat32_2DArray(t *testing.T) {
	// Test with valid input
	input := []interface{}{
		[]interface{}{float64(1.0), float32(2.0)},
		[]interface{}{float64(3.0), float32(4.0)},
	}
	expected := [][]float32{
		{1.0, 2.0},
		{3.0, 4.0},
	}
	result, err := ConvertToFloat32_2DArray(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test with invalid input
	input = []interface{}{
		[]interface{}{float64(1.0), "invalid"},
	}
	_, err = ConvertToFloat32_2DArray(input)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
