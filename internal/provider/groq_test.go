package provider

import "testing"

func TestGroqSendGemmaFamily(t *testing.T) {
	c, err := NewGroq("")
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Gemma 2 9B IT")
}

func TestGroqSendLlama3Family(t *testing.T) {
	c, err := NewGroq("")
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Llama 3.3 70b Versatile")
	sendHelper(t, c, "Llama 3.1 8b Instant")
	sendHelper(t, c, "Llama Guard 3 8B")
	sendHelper(t, c, "Llama3 70b 8192")
	sendHelper(t, c, "Llama3 8b 8192")
	sendHelper(t, c, "Llama 3.2 1b Preview")
	sendHelper(t, c, "Llama 3.2 3b Preview")
}

func TestGroqSendMixtralFamily(t *testing.T) {
	c, err := NewGroq("")
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "Mixtral 8x7b 32768")
}

func TestGroqSendR1Family(t *testing.T) {
	c, err := NewGroq("")
	if err != nil {
		t.Fatal(err)
	}

	sendHelper(t, c, "DeepSeek R1 Distill Llama 70B")
}

func sendHelper(t *testing.T, c Provider, modelName string) {
	m := c.GetLLMModels(modelName)
	if len(m) == 0 {
		t.Fatal(modelName + " should not be empty")
	}

	res, err := c.Send("Hey, whats your name?", m[0])
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Completion) == 0 {
		t.Error("completion should not be empty")
	}

	t.Logf("Model respose: %s", res.Completion)
}
