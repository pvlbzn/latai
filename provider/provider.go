package provider

import (
	"time"

	"github.com/pvlbzn/genlat/prompt"
)

// Provider is a core interface for each provider implementation
// to satisfy. Each service, such as AWS Bedrock or OpenAI, should
// provide implementation of Provider.
type Provider interface {
	// GetLLMModels returns a list of LLM models from memory.
	GetLLMModels(filter string) ([]*Model, error)

	// GetAllModels returns a list of all available models at provider.
	GetAllModels(filter string) ([]*Model, error)

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

	ModelProviderBedrock ModelProvider = "Bedrock"
	ModelProviderOpenAI  ModelProvider = "Open AI"

	ModelVendorOpenAI      ModelVendor = "Open AI"
	ModelVendorAmazon      ModelVendor = "Amazon"
	ModelVendorStabilityAI ModelVendor = "Stability AI"
	ModelVendorAI21Labs    ModelVendor = "AI21 Labs"
	ModelVendorAnthropic   ModelVendor = "Anthropic"
	ModelVendorCohere      ModelVendor = "Cohere"
	ModelVendorMeta        ModelVendor = "Meta"
	ModelVendorMistralAI   ModelVendor = "Mistral AI"
)
