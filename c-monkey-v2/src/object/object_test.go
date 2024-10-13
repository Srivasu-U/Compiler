package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}
	boolean1 := &Boolean{Value: true}
	boolean2 := &Boolean{Value: true}
	number1 := &Integer{Value: 100}
	number2 := &Integer{Value: 100}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if boolean1.HashKey() != boolean2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}

	if number1.HashKey() != number2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("string wtih different content have same hash keys")
	}
}
