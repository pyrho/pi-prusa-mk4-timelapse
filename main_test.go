package main

import (
	"testing"
)

func TestParseGcode(t *testing.T) {
	cmd := parseGcode("// status:print_start")
	if cmd != COMMAND_PRINT_START {
		t.Fatalf("Wrong command %v", cmd)
	}

	cmd = parseGcode("// status:print_stop")
	if cmd != COMMAND_PRINT_STOP {
		t.Fatalf("Wrong command %v", cmd)
	}

	cmd = parseGcode("// action:capture")
	if cmd != COMMAND_CAPTURE {
		t.Fatalf("Wrong command %v", cmd)
	}

	cmd = parseGcode("// random:shit")
	if cmd != COMMAND_UNHANDLED {
		t.Fatalf("Wrong command %v", cmd)
	}
}
