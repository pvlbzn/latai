package prompt

import (
	"testing"
)

func TestGetPrompts(t *testing.T) {
	res, err := GetPrompts()
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Fatal("wrong number of prompts loaded")
	}

	p1, p2, p3 := len(res[0].Content), len(res[1].Content), len(res[2].Content)

	if !(p1 > 0 && p2 > 0 && p3 > 0) {
		t.Fatalf("prompts shouldn't be empty, got `p1=%d`, `p2=%d`, `p3=%d`", p1, p2, p3)
	}
}
