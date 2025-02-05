package provider

import "testing"

func TestGroqSendGemmaFamily(t *testing.T) {
	groqSendHelper(t, "Gemma 2 9B IT")
}

func TestGroqSendLlama3Family(t *testing.T) {
	groqSendHelper(t, "Llama 3.3 70b Versatile")
	groqSendHelper(t, "Llama 3.1 8b Instant")
	groqSendHelper(t, "Llama Guard 3 8B")
	groqSendHelper(t, "Llama3 70b 8192")
	groqSendHelper(t, "Llama3 8b 8192")
	groqSendHelper(t, "Llama 3.2 1b Preview")
	groqSendHelper(t, "Llama 3.2 3b Preview")
}

func TestGroqSendMixtralFamily(t *testing.T) {
	groqSendHelper(t, "Mixtral 8x7b 32768")
}

func TestGroqSendR1Family(t *testing.T) {
	groqSendHelper(t, "DeepSeek R1 Distill Llama 70B")
}

func groqSendHelper(t *testing.T, modelName string) {
	c, err := NewGroq("")
	if err != nil {
		t.Fatal(err)
	}

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
