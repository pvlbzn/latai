package provider

import (
	"strings"
	"time"

	"github.com/pvlbzn/latai/prompt"
)

// Provider is a core interface for each provider implementation
// to satisfy. Each service, such as AWS Bedrock or OpenAI, should
// provide implementation of Provider.
type Provider interface {
	// GetLLMModels returns a list of LLM models from memory.
	GetLLMModels(filter string) []*Model

	// Measure measures a particular model and returns Metric back.
	Measure(model *Model, prompt *prompt.Prompt) (*Metric, error)

	// Send a message to LLM. Can be used stand alone, and is used
	// by Measure internally to make calls to gather metrics.
	Send(message string, to *Model) (*Response, error)
}

type Response struct {
	Completion string `json:"completion"`
}

// Metric wraps model data and provides Latency extra field.
type Metric struct {
	Model    *Model
	Latency  time.Duration
	Response *Response
}

// Model holds key characteristics of a particular model instance.
type Model struct {
	// Model identification, oftentimes model's full name.
	ID string

	// Model human-readable name.
	Name string

	// Family represents family of model which often times defines its API.
	Family ModelFamily

	// Model service provider, that is name of a service model
	// is being served from such as AWS, OpenAI, Grok, etc.
	Provider ModelProvider

	// Model vendor, that is name of a company which built
	// the model itself such as Anthropic, Google, Amazon, etc.
	Vendor ModelVendor
}

type ModelFamily string
type ModelProvider string
type ModelVendor string

const (
	ModelFamilyTitan    ModelFamily = "Titan"
	ModelFamilyNova     ModelFamily = "Nova"
	ModelFamilyGPT      ModelFamily = "GPT"
	ModelFamilyClaude   ModelFamily = "Claude"
	ModelFamilyJurassic ModelFamily = "Jurassic"
	ModelFamilyJamba    ModelFamily = "Jamba"
	ModelFamilyCommand  ModelFamily = "Command"
	ModelFamilyCommandR ModelFamily = "Command R"
	ModelFamilyLlama3   ModelFamily = "Llama 3"
	ModelFamilyMistral  ModelFamily = "Mistral"
	ModelFamilyMixtral  ModelFamily = "Mixtral"
	ModelFamilyGemma    ModelFamily = "Gemma"
	ModelFamilyR1       ModelFamily = "R1"

	ModelProviderBedrock ModelProvider = "Bedrock"
	ModelProviderOpenAI  ModelProvider = "Open AI"
	ModelProviderGroq    ModelProvider = "Groq"

	ModelVendorOpenAI    ModelVendor = "Open AI"
	ModelVendorAmazon    ModelVendor = "Amazon"
	ModelVendorAI21Labs  ModelVendor = "AI21 Labs"
	ModelVendorAnthropic ModelVendor = "Anthropic"
	ModelVendorCohere    ModelVendor = "Cohere"
	ModelVendorMeta      ModelVendor = "Meta"
	ModelVendorMistralAI ModelVendor = "Mistral AI"
	ModelVendorGoogle    ModelVendor = "Google"
	ModelVendorDeepSeek  ModelVendor = "DeepSeek"
)

func measure(provider Provider, model *Model, prompt *prompt.Prompt) (*Metric, error) {
	start := time.Now()
	res, err := provider.Send(prompt.Content, model)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)

	return &Metric{
		Model:    model,
		Latency:  elapsed,
		Response: res,
	}, nil
}

// filterModels returns models which model name is a substring of filter
// string. If filter is empty string then all models returned (empty set
// is a subset of every set). If no models found then empty list returned.
func filterModels(models []Model, filter string) []*Model {
	// Pre-allocate list enough to hold all models to avoid reallocations.
	res := make([]*Model, 0, len(models))

	for _, model := range models {
		modelName, query := strings.ToLower(model.Name), strings.ToLower(filter)

		if strings.Contains(modelName, query) {
			modelCopy := model
			res = append(res, &modelCopy)
		}
	}

	return res
}
