package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/pvlbzn/genlat/prompt"
)

const (
	DefaultAWSRegion  string = "us-east-1"
	DefaultAWSProfile string = "default"

	BedrockVendorAmazon      string = "Amazon"
	BedrockVendorStabilityAI string = "Stability AI"
	BedrockVendorAI21Labs    string = "AI21 Labs"
	BedrockVendorAnthropic   string = "Anthropic"
	BedrockVendorCohere      string = "Cohere"
	BedrockVendorMeta        string = "Meta"
	BedrockVendorMistralAI   string = "Mistral AI"
)

type Bedrock struct {
	client  *bedrock.Client
	runtime *bedrockruntime.Client
	models  []Model
}

// NewBedrock creates a new AWS Bedrock client with provided region and profile.
// If you use default region and profile use DefaultAWSRegion
// and DefaultAWSProfile.
func NewBedrock(region string, profile string) (*Bedrock, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile))
	if err != nil {
		slog.Error("failed to load bedrock default configuration", "error", err)
		return nil, err
	}

	models := []Model{
		{ID: "amazon.titan-tg1-large", Name: "Titan Text Large", Provider: "AWS Bedrock", Vendor: "Amazon"},
		{ID: "amazon.titan-text-premier-v1:0", Name: "Titan Text G1 - Premier", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-pro-v1:0:300k", Name: "Nova Pro", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-pro-v1:0", Name: "Nova Pro", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-lite-v1:0:300k", Name: "Nova Lite", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-lite-v1:0", Name: "Nova Lite", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-canvas-v1:0", Name: "Nova Canvas", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-reel-v1:0", Name: "Nova Reel", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-micro-v1:0:128k", Name: "Nova Micro", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.nova-micro-v1:0", Name: "Nova Micro", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-g1-text-02", Name: "Titan Text Embeddings v2", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-text-lite-v1:0:4k", Name: "Titan Text G1 - Lite", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-text-lite-v1", Name: "Titan Text G1 - Lite", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-text-express-v1:0:8k", Name: "Titan Text G1 - Express", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-text-express-v1", Name: "Titan Text G1 - Express", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-text-v1:2:8k", Name: "Titan Embeddings G1 - Text", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-text-v1", Name: "Titan Embeddings G1 - Text", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-text-v2:0:8k", Name: "Titan Text Embeddings V2", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-text-v2:0", Name: "Titan Text Embeddings V2", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-image-v1:0", Name: "Titan Multimodal Embeddings G1", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "amazon.titan-embed-image-v1", Name: "Titan Multimodal Embeddings G1", Provider: "AWS Bedrock", Vendor: "Amazon"},
		//{ID: "stability.stable-diffusion-xl-v1:0", Name: "SDXL 1.0", Provider: "AWS Bedrock", Vendor: "Stability AI"},
		//{ID: "stability.stable-diffusion-xl-v1", Name: "SDXL 1.0", Provider: "AWS Bedrock", Vendor: "Stability AI"},
		//{ID: "ai21.j2-grande-instruct", Name: "J2 Grande Instruct", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.j2-jumbo-instruct", Name: "J2 Jumbo Instruct", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.j2-mid", Name: "Jurassic-2 Mid", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.j2-mid-v1", Name: "Jurassic-2 Mid", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.j2-ultra", Name: "Jurassic-2 Ultra", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.j2-ultra-v1:0:8k", Name: "Jurassic-2 Ultra", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.j2-ultra-v1", Name: "Jurassic-2 Ultra", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.jamba-instruct-v1:0", Name: "Jamba-Instruct", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.jamba-1-5-large-v1:0", Name: "Jamba 1.5 Large", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "ai21.jamba-1-5-mini-v1:0", Name: "Jamba 1.5 Mini", Provider: "AWS Bedrock", Vendor: "AI21 Labs"},
		//{ID: "cohere.command-text-v14:7:4k", Name: "Command", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.command-text-v14", Name: "Command", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.command-r-v1:0", Name: "Command R", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.command-r-plus-v1:0", Name: "Command R+", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.command-light-text-v14:7:4k", Name: "Command Light", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.command-light-text-v14", Name: "Command Light", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.embed-english-v3:0:512", Name: "Embed English", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.embed-english-v3", Name: "Embed English", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.embed-multilingual-v3:0:512", Name: "Embed Multilingual", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "cohere.embed-multilingual-v3", Name: "Embed Multilingual", Provider: "AWS Bedrock", Vendor: "Cohere"},
		//{ID: "meta.llama3-8b-instruct-v1:0", Name: "Llama 3 8B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-70b-instruct-v1:0", Name: "Llama 3 70B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-1-8b-instruct-v1:0", Name: "Llama 3.1 8B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-1-70b-instruct-v1:0", Name: "Llama 3.1 70B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-2-11b-instruct-v1:0", Name: "Llama 3.2 11B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-2-90b-instruct-v1:0", Name: "Llama 3.2 90B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-2-1b-instruct-v1:0", Name: "Llama 3.2 1B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-2-3b-instruct-v1:0", Name: "Llama 3.2 3B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "meta.llama3-3-70b-instruct-v1:0", Name: "Llama 3.3 70B Instruct", Provider: "AWS Bedrock", Vendor: "Meta"},
		//{ID: "mistral.mistral-7b-instruct-v0:2", Name: "Mistral 7B Instruct", Provider: "AWS Bedrock", Vendor: "Mistral AI"},
		//{ID: "mistral.mixtral-8x7b-instruct-v0:1", Name: "Mixtral 8x7B Instruct", Provider: "AWS Bedrock", Vendor: "Mistral AI"},
		//{ID: "mistral.mistral-large-2402-v1:0", Name: "Mistral Large (24.02)", Provider: "AWS Bedrock", Vendor: "Mistral AI"},
		//{ID: "mistral.mistral-small-2402-v1:0", Name: "Mistral Small (24.02)", Provider: "AWS Bedrock", Vendor: "Mistral AI"},

		// Claude family.
		//{ID: "anthropic.claude-instant-v1", Name: "Claude Instant v1", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "anthropic.claude-v2:1", Name: "Claude v2:1", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "anthropic.claude-v2", Name: "Claude v2", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "us.anthropic.claude-3-haiku-20240307-v1:0", Name: "Claude 3 Haiku", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "us.anthropic.claude-3-sonnet-20240229-v1:0", Name: "Claude 3 Sonnet", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "us.anthropic.claude-3-5-haiku-20241022-v1:0", Name: "Claude 3.5 Haiku", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "us.anthropic.claude-3-5-sonnet-20240620-v1:0", Name: "Claude 3.5 Sonnet v1", Provider: "AWS Bedrock", Vendor: "Anthropic"},
		//{ID: "us.anthropic.claude-3-5-sonnet-20241022-v2:0", Name: "Claude 3.5 Sonnet v2", Provider: "AWS Bedrock", Vendor: "Anthropic"},

	}

	return &Bedrock{
		client:  bedrock.NewFromConfig(cfg),
		runtime: bedrockruntime.NewFromConfig(cfg),
		models:  models,
	}, nil
}

// GetLLMModels returns LLM models only which name matches filter.
// Empty filter string returns all models unfiltered.
func (s *Bedrock) GetLLMModels(filter string) ([]*Model, error) {
	// Pre-allocate list enough to hold all models to avoid reallocations.
	res := make([]*Model, 0, len(s.models))

	for _, model := range s.models {
		modelName, query := strings.ToLower(model.Name), strings.ToLower(filter)

		if strings.Contains(modelName, query) {
			modelCopy := model
			res = append(res, &modelCopy)
		}
	}

	return res, nil
}

// GetAllModels from AWS Bedrock service. Effectively lists available models
// from API. Note that possibly not all models are LLM models. Some of them
// are embeddings, instruct models, etc.
func (s *Bedrock) GetAllModels(filter string) ([]*Model, error) {
	res, err := s.client.ListFoundationModels(
		context.TODO(),
		&bedrock.ListFoundationModelsInput{})
	if err != nil {
		slog.Error("couldn't get list of foundation models", "error", err.Error())
		return nil, err
	}

	if len(res.ModelSummaries) == 0 {
		m := "no foundation models found"
		slog.Error(m)
		return nil, fmt.Errorf(m)
	}

	models := make([]*Model, 0, len(res.ModelSummaries))
	for _, summary := range res.ModelSummaries {
		modelName, query := strings.ToLower(*summary.ModelName), strings.ToLower(filter)

		if strings.Contains(modelName, query) {
			models = append(
				models,
				&Model{
					ID:       *summary.ModelId,
					Name:     *summary.ModelName,
					Provider: "AWS Bedrock",
					Vendor:   *summary.ProviderName,
				})
		}
	}

	return models, nil
}

func (s *Bedrock) Measure(model *Model, prompt *prompt.Prompt) (*Metric, error) {
	panic("not implemented")
}

// Send message.
func (s *Bedrock) Send(message string, to *Model) (*Response, error) {
	// Internally Send is a routing function which delegates actual
	// computation to an appropriate vendor handler.
	switch to.Vendor {
	case BedrockVendorAmazon:
		return s.runTitanFamilyInference(message, to)

	case BedrockVendorStabilityAI:
		return s.runInferenceStabilityAI(message, to)

	case BedrockVendorAI21Labs:
		return s.runInferenceAI21Labs(message, to)

	case BedrockVendorAnthropic:
		return s.runClaudeFamilyInference(message, to)

	case BedrockVendorCohere:
		return s.runInferenceCohere(message, to)

	case BedrockVendorMeta:
		return s.runInferenceMeta(message, to)

	case BedrockVendorMistralAI:
		return s.runInferenceMistralAI(message, to)

	default:
		return nil, fmt.Errorf("unsupported model vendor: %s", to.Vendor)
	}
}

type modelParserFn func(output *bedrockruntime.InvokeModelOutput) (*Response, error)

type titanRequest struct {
	InputText            string               `json:"inputText"`
	TextGenerationConfig textGenerationConfig `json:"textGenerationConfig"`
}

type textGenerationConfig struct {
	MaxTokenCount int      `json:"maxTokenCount"`
	Temperature   float32  `json:"temperature"`
	TopP          float32  `json:"topP"`
	StopSequences []string `json:"stopSequences"`
}

type titanResponse struct {
	Results []struct {
		OutputText string `json:"outputText"`
	} `json:"results"`
}

func (s *Bedrock) runTitanFamilyInference(message string, to *Model) (*Response, error) {
	data := &titanRequest{
		InputText: message,
		TextGenerationConfig: textGenerationConfig{
			MaxTokenCount: 1024,
			Temperature:   0.1,
			TopP:          0.5,
			StopSequences: []string{},
		},
	}

	parser := func(output *bedrockruntime.InvokeModelOutput) (*Response, error) {
		var res titanResponse
		err := json.Unmarshal(output.Body, &res)
		if err != nil {
			slog.Debug("failed to unmarshal response", "error", err.Error(), "model", *to, "data", data)
			return nil, err
		}

		return &Response{
			Completion: res.Results[0].OutputText,
		}, nil
	}

	return s.runInference(data, to, parser)
}

func (s *Bedrock) runInferenceStabilityAI(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

func (s *Bedrock) runInferenceAI21Labs(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

type claudeRequest struct {
	Messages         []claudeMessage `json:"messages"`
	MaxTokens        int             `json:"max_tokens"`
	Temperature      float64         `json:"temperature"`
	TopP             float64         `json:"top_p"`
	AnthropicVersion string          `json:"anthropic_version"`
}

type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (s *Bedrock) runClaudeFamilyInference(message string, to *Model) (*Response, error) {
	data := claudeRequest{
		Messages: []claudeMessage{
			{Role: "user", Content: message},
		},
		MaxTokens:        1024,
		Temperature:      0.5,
		TopP:             0.5,
		AnthropicVersion: "bedrock-2023-05-31",
	}

	parser := func(output *bedrockruntime.InvokeModelOutput) (*Response, error) {
		var res claudeResponse
		err := json.Unmarshal(output.Body, &res)
		if err != nil {
			slog.Debug("failed to unmarshal response", "error", err.Error(), "model", *to, "data", data)
			return nil, err
		}

		return &Response{
			Completion: res.Content[0].Text,
		}, nil
	}

	return s.runInference(data, to, parser)
}

func (s *Bedrock) runInferenceCohere(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

func (s *Bedrock) runInferenceMeta(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

func (s *Bedrock) runInferenceMistralAI(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

func (s *Bedrock) runInference(data any, to *Model, withParser modelParserFn) (*Response, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		slog.Debug("failed to marshal model data", "error", err.Error(), "data", data)
		return nil, err
	}

	out, err := s.runtime.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(to.ID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        dataBytes,
	})
	if err != nil {
		slog.Debug("failed to invoke model", "error", err.Error(), "model", *to, "data", data)
		return nil, err
	}

	return withParser(out)
}
