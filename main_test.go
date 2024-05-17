package main

import (
    "testing"
)
// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestParseGcode(t *testing.T) {
    // name := "Gladys"
    // want := regexp.MustCompile(`\b`+name+`\b`)
    msg := parseGcode("Gladys")
    t.Log(msg)
}
