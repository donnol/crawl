package main

import "testing"

func TestPhantom(t *testing.T) {
	testCase := []string{
		"http://www.taobao.com",
	}

	for _, c := range testCase {
		content, err := phantom(c)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s\n", content)
	}
}
