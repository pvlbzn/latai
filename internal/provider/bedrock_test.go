package provider

import (
	"testing"
)

// Mistral Family
func TestBedrockSendMistralFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Mistral Large (24.02)")
	sendHelper(t, c, "Mistral Small (24.02)")
}

// Meta Family
func TestBedrockSendLlama3Family(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Llama 3 8B Instruct")
	sendHelper(t, c, "Llama 3 70B Instruct")
}

// Command Family
func TestBedrockSendCommandFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Command")
	sendHelper(t, c, "Command Light")
}

// Command R Family
func TestBedrockSendCommandRFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Command R")
	sendHelper(t, c, "Command R+")
}

// Jamba Family
func TestBedrockSendJambaFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Jamba 1.5 Large")
	sendHelper(t, c, "Jamba 1.5 Mini")
}

// Jurassic Family
func TestBedrockSendJurassicFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Jurassic-2 Mid")
	sendHelper(t, c, "Jurassic-2 Mid")
	sendHelper(t, c, "Jurassic-2 Ultra")
}

// Nova Family
func TestBedrockSendNovaFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Nova Pro")
	sendHelper(t, c, "Nova Lite")
	sendHelper(t, c, "Nova Micro")
}

// Titan Family
func TestBedrockSendTitanFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Titan Text Large")
	sendHelper(t, c, "Titan Text G1 - Premier")
	sendHelper(t, c, "Titan Text G1 - Lite")
	sendHelper(t, c, "Titan Text G1 - Express")
}

// Claude Family
func TestBedrockSendClaudeFamily(t *testing.T) {
	c, err := NewBedrock(DefaultAWSRegion, DefaultAWSProfile)
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Claude Instant v1")
	sendHelper(t, c, "Claude v2:1")
	sendHelper(t, c, "Claude v2")
	sendHelper(t, c, "Claude 3 Haiku")
	sendHelper(t, c, "Claude 3 Sonnet")
	sendHelper(t, c, "Claude 3.5 Haiku")
	sendHelper(t, c, "Claude 3.5 Sonnet v1")
	sendHelper(t, c, "Claude 3.5 Sonnet v2")
}
