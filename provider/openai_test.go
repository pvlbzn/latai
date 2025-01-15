package provider

import (
	"strings"
	"testing"
)

func TestOpenAIid2Name(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{
			in:       "gpt-4o-realtime-preview",
			expected: "Gpt 4o Realtime Preview",
		},
		{
			in:       "o1-mini",
			expected: "O1 Mini",
		},
	}

	client, err := NewOpenAI("")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		res := client.id2Name(tt.in)
		if res != tt.expected {
			t.Errorf("id2Name(%s) got %s, want %s", tt.in, res, tt.expected)
		}
	}
}

// TestGetModels executes `GetModels` and verifies that models returned
// correctly. It verifies just a few models out of the list because
// the list gets updates with some frequency.
func TestOpenAIGetModels(t *testing.T) {
	expectedModels := []string{
		"gpt-4o-realtime-preview-2024-12-17",
		"gpt-4o-mini-audio-preview",
		"dall-e-3",
		"o1-preview",
		"gpt-4-turbo-2024-04-09",
	}

	client, err := NewOpenAI("")
	if err != nil {
		t.Fatal(err)
	}

	models, err := client.GetAllModels("")
	if err != nil {
		t.Fatal(err)
	}

	// Build a lookup table so all the future checks will be
	// O(1) instead of O(n)
	lookup := make(map[string]string)
	for _, m := range models {
		lookup[m.ID] = m.Name
	}

	for _, m := range expectedModels {
		_, found := lookup[m]
		if !found {
			t.Errorf("model `%s` not found", m)
		}
	}
}

func TestOpenAIGetLLMModelsSearch(t *testing.T) {
	client, err := NewOpenAI("")
	if err != nil {
		t.Fatal(err)
	}

	models, err := client.GetLLMModels("o1 mini")
	if err != nil {
		t.Fatal(err)
	}

	if len(models) != 2 {
		t.Errorf("got %d models, expected 2", len(models))
	}

	if models[0].ID != "o1-mini" {
		t.Errorf("got %s, expected o1-mini", models[0].ID)
	}

	if models[1].ID != "o1-mini-2024-09-12" {
		t.Errorf("got %s, expected o1-mini-2024-09-12", models[1].ID)
	}
}

func TestOpenAISend(t *testing.T) {
	client, err := NewOpenAI("")
	if err != nil {
		t.Fatal(err)
	}

	models, err := client.GetLLMModels("o1 mini")
	if err != nil {
		t.Fatal(err)
	}

	prompt := "Reply a single word 'Assistant'"
	expected := "Assistant"

	res, err := client.Send(prompt, models[0])
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(res) != expected {
		t.Errorf("got %s, want %s", res, expected)
	}
}
