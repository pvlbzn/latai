package provider

import (
	"strings"
	"testing"
)

// TestGetModels executes `GetModels` and verifies that models returned
// correctly. It verifies just a few models out of the list because
// the list gets updates with some frequency.
func TestBedrockGetModels(t *testing.T) {
	type model struct {
		id       string
		model    string
		vendor   string
		provider string
	}
	tests := []model{
		{
			id:       "mistral.mixtral-8x7b-instruct-v0:1",
			model:    "Mixtral 8x7B Instruct",
			vendor:   "Mistral AI",
			provider: "AWS Bedrock",
		},
		{
			id:       "amazon.titan-image-generator-v2:0",
			model:    "Titan Image Generator G1 v2",
			vendor:   "Amazon,",
			provider: "AWS Bedrock",
		},
		{
			id:       "ai21.j2-ultra-v1",
			model:    "Jurassic-2 Ultra",
			vendor:   "AI21 Labs",
			provider: "AWS Bedrock",
		},
		{
			id:       "anthropic.claude-3-sonnet-20240229-v1:0:200k",
			model:    "Claude 3 Sonnet",
			vendor:   "Anthropic,",
			provider: "AWS Bedrock",
		},
		{
			id:       "cohere.command-r-plus-v1:0",
			model:    "Command R+",
			vendor:   "Cohere,",
			provider: "AWS Bedrock",
		},
	}

	client, err := NewBedrock(DefaultBedrockRegion)
	if err != nil {
		t.Fatal(err)
	}

	models, err := client.GetModels()
	if err != nil {
		t.Fatal(err)
	}

	// Build a lookup table so all the future checks will be
	// O(1) instead of O(n) resulting in O(n) + O(n) complexity
	// instead of O(n^2)
	lookup := make(map[string]model)
	for _, m := range models {
		lookup[m.ID] = model{
			id:       m.ID,
			model:    m.Name,
			vendor:   m.Vendor,
			provider: m.Provider,
		}
	}

	for _, m := range tests {
		_, found := lookup[m.id]
		if !found {
			t.Errorf("model `%s` not found", m)
		}
	}
}

func TestBedrockSend(t *testing.T) {
	client, err := NewBedrock(DefaultBedrockRegion)
	if err != nil {
		t.Fatal(err)
	}

	prompt :=
		"What is your name? Reply in a single word, without punctuation or anything else."
	expected :=
		"Claude"

	res, err := client.Send(prompt)
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(res) != expected {
		t.Errorf("got %s, want %s", res, expected)
	}
}
