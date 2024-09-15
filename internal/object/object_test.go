package object

import (
	"github.com/VannRR/golox/internal/value"
	"testing"
)

func Test_ObjString_IsType(t *testing.T) {
	objString := ObjString("foo")
	otherValue := value.BoolVal(true)

	if !objString.IsType(objString) {
		t.Errorf("Expected IsType to return true for the same ObjString type, but got false")
	}

	if objString.IsType(otherValue) {
		t.Errorf("Expected IsType to return false for different types, but got true")
	}
}

func Test_ObjString_IsEqual(t *testing.T) {
	foo := ObjString("foo")
	bar := ObjString("bar")
	otherValue := value.NumberVal(1)

	if foo.IsEqual(bar) {
		t.Errorf("Expected IsEqual to return false for ObjString 'foo' == 'bar', but got true")
	}

	if !foo.IsEqual(foo) {
		t.Errorf("Expected IsEqual to return true for ObjString 'foo' == 'foo', but got false")
	}

	if foo.IsEqual(otherValue) {
		t.Errorf("Expected IsEqual to return false for different types, but got true")
	}
}

func Test_ObjString_IsFalsey(t *testing.T) {
	str := ObjString("wow")

	if str.IsFalsey() {
		t.Errorf("Expected IsFalsey to return false for ObjString, but got true")
	}
}

func Test_ObjString_Stringify(t *testing.T) {
	foo := ObjString("foo")
	bar := ObjString("bar")

	expectedFooString := "foo"
	actualFooString := foo.String()

	if actualFooString != expectedFooString {
		t.Errorf("Expected Stringify to return \"%s\" for ObjString 'foo', but got \"%s\"", expectedFooString, actualFooString)
	}

	expectedBarString := "bar"
	actualBarString := bar.String()

	if actualBarString != expectedBarString {
		t.Errorf("Expected Stringify to return \"%s\" for ObjString 'bar', but got \"%s\"", expectedBarString, actualBarString)
	}
}

func Test_ObjFuntion_IsType(t *testing.T) {
	objFunction := NewFunction()
	otherValue := value.BoolVal(true)

	if !objFunction.IsType(objFunction) {
		t.Errorf("Expected IsType to return true for the same objFunction type, but got false")
	}

	if objFunction.IsType(otherValue) {
		t.Errorf("Expected IsType to return false for different types, but got true")
	}
}

func Test_ObjFuntion_IsEqual(t *testing.T) {
	foo := NewFunction()
	foo.name = "foo"
	bar := NewFunction()
	bar.name = "bar"
	otherValue := value.NumberVal(1)

	if foo.IsEqual(*bar) {
		t.Errorf("Expected IsEqual to return false for ObjFunction '<fn foo>' == '<fn bar>', but got true")
	}

	if foo.IsEqual(*foo) {
		t.Errorf("Expected IsEqual to return false for ObjFunction '<fn foo>' == '<fn foo>', but got true")
	}

	if foo.IsEqual(otherValue) {
		t.Errorf("Expected IsEqual to return false for different types, but got true")
	}
}

func Test_ObjFunction_IsFalsey(t *testing.T) {
	objFunction := NewFunction()

	if objFunction.IsFalsey() {
		t.Errorf("Expected IsFalsey to return false for ObjFunction, but got true")
	}
}

func Test_ObjFunction_Stringify(t *testing.T) {
	foo := NewFunction()
	foo.name = "foo"
	bar := NewFunction()
	bar.name = "bar"

	expectedFooString := "<fn foo>"
	actualFooString := foo.String()

	if actualFooString != expectedFooString {
		t.Errorf("Expected Stringify to return \"%s\" for ObjFunction 'foo', but got \"%s\"", expectedFooString, actualFooString)
	}

	expectedBarString := "<fn bar>"
	actualBarString := bar.String()

	if actualBarString != expectedBarString {
		t.Errorf("Expected Stringify to return \"%s\" for ObjFunction 'bar', but got \"%s\"", expectedBarString, actualBarString)
	}
}
