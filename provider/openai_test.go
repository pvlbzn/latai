package provider

import "testing"

func TestOpenAISendGPTFamily(t *testing.T) {
	c, err := NewOpenAI("")
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "GPT 4 1106 Preview")
	sendHelper(t, c, "GPT 3.5 Turbo")
	sendHelper(t, c, "GPT 3.5 Turbo 0125")
	sendHelper(t, c, "GPT 3.5 Turbo 16k")
	sendHelper(t, c, "GPT 4")
	sendHelper(t, c, "GPT 4 0613")
	sendHelper(t, c, "ChatGPT 4o Latest")
	sendHelper(t, c, "GPT 4o 2024 08 06")
	sendHelper(t, c, "GPT 4o")
	sendHelper(t, c, "GPT 3.5 Turbo 1106")
	sendHelper(t, c, "GPT 4 Turbo 2024 04 09")
	sendHelper(t, c, "GPT 4 Turbo")
	sendHelper(t, c, "GPT 4 Turbo Preview")
	sendHelper(t, c, "GPT 4o 2024 05 13")
	sendHelper(t, c, "GPT 4o 2024 11 20")
	sendHelper(t, c, "GPT 4o Mini 2024 07 18")
	sendHelper(t, c, "GPT 4o Mini")
	sendHelper(t, c, "GPT 4 0125 Preview")

	// These tests are disabled due to their price, and likelihood
	// of being correct if all the above runs successfully since
	// API is identical.
	/*
		sendHelper(t, c, "O1")
		sendHelper(t, c, "O1 Preview 2024 09 12")
		sendHelper(t, c, "O1 Preview")
		sendHelper(t, c, "O1 Mini")
		sendHelper(t, c, "O1 Mini 2024 09 12")
		sendHelper(t, c, "O1 2024 12 17")
	*/
}
