package provider

import (
	"time"

	"github.com/pvlbzn/genlat/prompt"
)

// Provider is a core interface for each provider implementation
// to satisfy. Each service, such as AWS Bedrock or OpenAI, should
// provide implementation of Provider.
type Provider interface {
	// GetModels returns list of available models at provider.
	GetModels(filter string) ([]*Model, error)

	// Measure measures a particular model and returns Metric back.
	Measure(model *Model, prompt *prompt.Prompt) (*Metric, error)

	// Send a message to LLM. Can be used stand alone, and is used
	// by Measure internally to make calls to gather metrics.
	Send(message string, to *Model) (string, error)
}

// Model holds key characteristics of a particular model instance.
type Model struct {
	// Model identification, oftentimes model's full name.
	ID string

	// Model human-readable name.
	Name string

	// Model service provider, that is name of a service model
	// is being served from such as AWS, OpenAI, Grok, etc.
	Provider string

	// Model vendor, that is name of a company which built
	// the model itself such as Anthropic, Google, Amazon, etc.
	Vendor string
}

// Metric wraps model data and provides Latency extra field.
type Metric struct {
	Model    *Model
	Latency  time.Duration
	Response string
}
