package value

import (
	"testing"
)

func Test_NilVal_IsType(t *testing.T) {
	nilValue := NilVal{}
	otherValue := NumberVal(1)

	if !nilValue.IsType(nilValue) {
		t.Errorf("Expected IsType to return true for the same NilVal type, but got false")
	}

	if nilValue.IsType(otherValue) {
		t.Errorf("Expected IsType to return false for different types, but got true")
	}
}

func Test_NilVal_IsEqual(t *testing.T) {
	nilValue := NilVal{}
	otherValue := NumberVal(1)

	if !nilValue.IsEqual(nilValue) {
		t.Errorf("Expected IsEqual to return true for the same NilVal, but got false")
	}

	if nilValue.IsEqual(otherValue) {
		t.Errorf("Expected IsEqual to return false for different types, but got true")
	}
}

func Test_NilVal_IsFalsey(t *testing.T) {
	nilValue := NilVal{}

	if !nilValue.IsFalsey() {
		t.Errorf("Expected IsFalsey to return true for NilVal, but got false")
	}
}

func Test_NilVal_Stringify(t *testing.T) {
	nilValue := NilVal{}

	expectedString := "nil"
	actualString := nilValue.Stringify()

	if actualString != expectedString {
		t.Errorf("Expected Stringify to return \"%s\", but got \"%s\"", expectedString, actualString)
	}
}

func Test_BoolVal_IsType(t *testing.T) {
	boolValue := BoolVal(true)
	otherValue := NumberVal(1)

	if !boolValue.IsType(boolValue) {
		t.Errorf("Expected IsType to return true for the same BoolVal type, but got false")
	}

	if boolValue.IsType(otherValue) {
		t.Errorf("Expected IsType to return false for different types, but got true")
	}
}

func Test_BoolVal_IsEqual(t *testing.T) {
	trueValue := BoolVal(true)
	falseValue := BoolVal(false)
	otherValue := NumberVal(123)

	if !trueValue.IsEqual(trueValue) {
		t.Errorf("Expected IsEqual to return true for the same BoolVal, but got false")
	}

	if trueValue.IsEqual(falseValue) {
		t.Errorf("Expected IsEqual to return false for different bool values, but got true")
	}

	if trueValue.IsEqual(otherValue) {
		t.Errorf("Expected IsEqual to return false for different types, but got true")
	}
}

func Test_BoolVal_IsFalsey(t *testing.T) {
	trueValue := BoolVal(true)
	falseValue := BoolVal(false)

	if trueValue.IsFalsey() {
		t.Errorf("Expected IsFalsey to return false for true BoolVal, but got true")
	}

	if !falseValue.IsFalsey() {
		t.Errorf("Expected IsFalsey to return true for false BoolVal, but got false")
	}
}

func Test_BoolVal_Stringify(t *testing.T) {
	trueValue := BoolVal(true)
	falseValue := BoolVal(false)

	expectedTrueString := "true"
	actualTrueString := trueValue.Stringify()

	if actualTrueString != expectedTrueString {
		t.Errorf("Expected Stringify to return \"%s\" for true BoolVal, but got \"%s\"", expectedTrueString, actualTrueString)
	}

	expectedFalseString := "false"
	actualFalseString := falseValue.Stringify()

	if actualFalseString != expectedFalseString {
		t.Errorf("Expected Stringify to return \"%s\" for false BoolVal, but got \"%s\"", expectedFalseString, actualFalseString)
	}
}

func Test_NumberVal_IsType(t *testing.T) {
	numberValue := NumberVal(1)
	otherValue := BoolVal(true)

	if !numberValue.IsType(numberValue) {
		t.Errorf("Expected IsType to return true for the same NumberVal type, but got false")
	}

	if numberValue.IsType(otherValue) {
		t.Errorf("Expected IsType to return false for different types, but got true")
	}
}

func Test_NumberVal_IsEqual(t *testing.T) {
	one := NumberVal(1)
	two := NumberVal(2)
	otherValue := BoolVal(true)

	if one.IsEqual(two) {
		t.Errorf("Expected IsEqual to return false for NumberVal 1 == 2, but got true")
	}

	if !one.IsEqual(one) {
		t.Errorf("Expected IsEqual to return true for NumberVal 1 == 1, but got false")
	}

	if one.IsEqual(otherValue) {
		t.Errorf("Expected IsEqual to return false for different types, but got true")
	}
}

func Test_NumberVal_IsFalsey(t *testing.T) {
	number := NumberVal(1)

	if number.IsFalsey() {
		t.Errorf("Expected IsFalsey to return false for NumberVal, but got true")
	}
}

func Test_NumberVal_Stringify(t *testing.T) {
	one := NumberVal(1)
	two := NumberVal(2)

	expectedOneString := "1"
	actualOneString := one.Stringify()

	if actualOneString != expectedOneString {
		t.Errorf("Expected Stringify to return \"%s\" for NumberVal 1, but got \"%s\"", expectedOneString, actualOneString)
	}

	expectedTwoString := "2"
	actualTwoString := two.Stringify()

	if actualTwoString != expectedTwoString {
		t.Errorf("Expected Stringify to return \"%s\" for NumberVal 2, but got \"%s\"", expectedTwoString, actualTwoString)
	}
}
