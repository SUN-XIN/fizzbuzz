package main

import "testing"

func TestProcess(t *testing.T) {
	// Test 1
	testClientRequest := clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    3,
		Int2:    5,
		Limit:   16,
	}

	responseOK := "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16"

	res := processRequest(&testClientRequest)
	if res != responseOK {
		t.Errorf("Test1 failed, expect\n%s\nbut get\n%s\n", responseOK, res)
	}

	// Test 2
	testClientRequest = clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    3,
		Int2:    3,
		Limit:   10,
	}

	responseOK = "1,2,fizzbuzz,4,5,fizzbuzz,7,8,fizzbuzz,10"
	res = processRequest(&testClientRequest)
	if res != responseOK {
		t.Errorf("Test2 failed, expect\n%s\nbut get\n%s\n", responseOK, res)
	}
}

func TestProcessBis(t *testing.T) {
	// Test 1
	testClientRequest := clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    3,
		Int2:    5,
		Limit:   16,
	}

	responseOK := "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,fizz,buzz,16"

	res := processRequestBis(&testClientRequest)
	if res != responseOK {
		t.Errorf("Test1 failed, expect\n%s\nbut get\n%s\n", responseOK, res)
	}

	// Test 2
	testClientRequest = clientRequest{
		String1: "fizz",
		String2: "buzz",
		Int1:    3,
		Int2:    3,
		Limit:   10,
	}

	responseOK = "1,2,fizzbuzz,fizz,buzz,4,5,fizzbuzz,fizz,buzz,7,8,fizzbuzz,fizz,buzz,10"

	res = processRequestBis(&testClientRequest)
	if res != responseOK {
		t.Errorf("Test1 failed, expect\n%s\nbut get\n%s\n", responseOK, res)
	}
}
