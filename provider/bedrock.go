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

		// Jamba family.
		{ID: "ai21.jamba-1-5-large-v1:0", Name: "Jamba 1.5 Large", Provider: ModelProviderBedrock, Vendor: ModelVendorAI21Labs, Family: ModelFamilyJamba},
		{ID: "ai21.jamba-1-5-mini-v1:0", Name: "Jamba 1.5 Mini", Provider: ModelProviderBedrock, Vendor: ModelVendorAI21Labs, Family: ModelFamilyJamba},

		// Jurassic family.
		{ID: "ai21.j2-mid", Name: "Jurassic-2 Mid", Provider: ModelProviderBedrock, Vendor: ModelVendorAI21Labs},
		{ID: "ai21.j2-mid-v1", Name: "Jurassic-2 Mid", Provider: ModelProviderBedrock, Vendor: ModelVendorAI21Labs},
		{ID: "ai21.j2-ultra", Name: "Jurassic-2 Ultra", Provider: ModelProviderBedrock, Vendor: ModelVendorAI21Labs},
		{ID: "ai21.j2-ultra-v1", Name: "Jurassic-2 Ultra", Provider: ModelProviderBedrock, Vendor: ModelVendorAI21Labs},

		// Nova family.
		{ID: "amazon.nova-pro-v1:0", Name: "Nova Pro", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyNova},
		{ID: "amazon.nova-lite-v1:0", Name: "Nova Lite", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyNova},
		{ID: "amazon.nova-micro-v1:0", Name: "Nova Micro", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyNova},

		//Titan family.
		{ID: "amazon.titan-tg1-large", Name: "Titan Text Large", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyTitan},
		{ID: "amazon.titan-text-premier-v1:0", Name: "Titan Text G1 - Premier", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyTitan},
		{ID: "amazon.titan-text-lite-v1", Name: "Titan Text G1 - Lite", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyTitan},
		{ID: "amazon.titan-text-express-v1", Name: "Titan Text G1 - Express", Provider: ModelProviderBedrock, Vendor: ModelVendorAmazon, Family: ModelFamilyTitan},

		// Claude family.
		{ID: "anthropic.claude-instant-v1", Name: "Claude Instant v1", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "anthropic.claude-v2:1", Name: "Claude v2:1", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "anthropic.claude-v2", Name: "Claude v2", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "us.anthropic.claude-3-haiku-20240307-v1:0", Name: "Claude 3 Haiku", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "us.anthropic.claude-3-sonnet-20240229-v1:0", Name: "Claude 3 Sonnet", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "us.anthropic.claude-3-5-haiku-20241022-v1:0", Name: "Claude 3.5 Haiku", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "us.anthropic.claude-3-5-sonnet-20240620-v1:0", Name: "Claude 3.5 Sonnet v1", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
		{ID: "us.anthropic.claude-3-5-sonnet-20241022-v2:0", Name: "Claude 3.5 Sonnet v2", Provider: ModelProviderBedrock, Vendor: ModelVendorAnthropic},
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
					Vendor:   ModelVendor(*summary.ProviderName),
				})
		}
	}

	return models, nil
}

func (s *Bedrock) Measure(model *Model, prompt *prompt.Prompt) (*Metric, error) {
	panic("not implemented")
}

// Send message.
func (s *Bedrock) Send(message string, model *Model) (*Response, error) {
	// Internally Send is a routing function which delegates actual
	// computation to an appropriate vendor handler.
	switch model.Vendor {
	case ModelVendorAmazon:
		switch model.Family {
		case ModelFamilyNova:
			return s.runBedrockInferenceNovaFamily(message, model)
		case ModelFamilyTitan:
			return s.runBedrockInferenceTitanFamily(message, model)
		default:
			return nil, fmt.Errorf("unsupported model family %s", model.Family)
		}

	case ModelVendorAI21Labs:
		switch model.Family {
		case ModelFamilyJurassic:
			return s.runBedrockInferenceJurassicFamily(message, model)
		case ModelFamilyJamba:
			return s.runBedrockInferenceJambaFamily(message, model)
		default:
			return nil, fmt.Errorf("unsupported model family %s", model.Family)
		}

	case ModelVendorAnthropic:
		return s.runBedrockInferenceClaudeFamily(message, model)

	case ModelVendorCohere:
		return s.runBedrockInferenceCohere(message, model)

	case ModelVendorMeta:
		return s.runBedrockInferenceMeta(message, model)

	case ModelVendorMistralAI:
		return s.runBedrockInferenceMistralAI(message, model)

	default:
		return nil, fmt.Errorf("unsupported model vendor: %s", model.Vendor)
	}
}

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

func (s *Bedrock) runBedrockInferenceTitanFamily(message string, model *Model) (*Response, error) {
	data := &titanRequest{
		InputText: message,
		TextGenerationConfig: textGenerationConfig{
			MaxTokenCount: 1024,
			Temperature:   0.1,
			TopP:          0.5,
			StopSequences: []string{},
		},
	}

	parser := func(res titanResponse) string {
		return res.Results[0].OutputText
	}

	return runBedrockInference(s, model, data, parser)
}

type novaRequest struct {
	Messages []novaMessage `json:"messages"`
}

type novaMessage struct {
	Role    string `json:"role"`
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

type novaResponse struct {
	Output struct {
		Message struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"message"`
	} `json:"output"`
}

func (s *Bedrock) runBedrockInferenceNovaFamily(message string, model *Model) (*Response, error) {
	data := &novaRequest{
		Messages: []novaMessage{
			{
				Role: "user",
				Content: []struct {
					Text string `json:"text"`
				}{{Text: message}},
			},
		},
	}

	parser := func(res novaResponse) string {
		return res.Output.Message.Content[0].Text
	}

	return runBedrockInference(s, model, data, parser)
}

type jurassicRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"maxTokens"`
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"topP"`
}

type jurassicResponse struct {
	Completions []struct {
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
	} `json:"completions"`
}

func (s *Bedrock) runBedrockInferenceJurassicFamily(message string, model *Model) (*Response, error) {
	data := jurassicRequest{
		Prompt:      message,
		MaxTokens:   1024,
		Temperature: 0.5,
		TopP:        0.5,
	}

	parser := func(res jurassicResponse) string {
		return res.Completions[0].Data.Text
	}

	return runBedrockInference(s, model, data, parser)
}

type jambaRequest struct {
	Messages []jambaMessage `json:"messages"`
}

type jambaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type jambaResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *Bedrock) runBedrockInferenceJambaFamily(message string, model *Model) (*Response, error) {
	data := &jambaRequest{
		Messages: []jambaMessage{
			{Role: "user", Content: message},
		},
	}

	parser := func(res jambaResponse) string {
		return res.Choices[0].Message.Content
	}

	return runBedrockInference(s, model, data, parser)
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

func (s *Bedrock) runBedrockInferenceClaudeFamily(message string, to *Model) (*Response, error) {
	data := claudeRequest{
		Messages: []claudeMessage{
			{Role: "user", Content: message},
		},
		MaxTokens:        1024,
		Temperature:      0.5,
		TopP:             0.5,
		AnthropicVersion: "bedrock-2023-05-31",
	}

	parser := func(in claudeResponse) string {
		return in.Content[0].Text
	}

	return runBedrockInference(s, to, data, parser)
}

func (s *Bedrock) runBedrockInferenceCohere(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

func (s *Bedrock) runBedrockInferenceMeta(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

func (s *Bedrock) runBedrockInferenceMistralAI(message string, to *Model) (*Response, error) {
	panic("not implemented")
}

// runBedrockInference is a helper function which wraps common Bedrock API operations.
// It receives Bedrock client, Model, and two generic types. The first generic
// is request object to a model, compliant to model's expected data. The second
// generic is model's output which is provided inside a parser. Parser unpacks
// model's response type into completion string.
func runBedrockInference[A, B any](bedrock *Bedrock, withModel *Model, withData A, withParser func(B) string) (*Response, error) {
	dataBytes, err := json.Marshal(withData)
	if err != nil {
		slog.Debug("failed to marshal model data", "error", err.Error(), "data", withData)
		return nil, err
	}

	out, err := bedrock.runtime.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(withModel.ID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        dataBytes,
	})
	if err != nil {
		slog.Debug("failed to invoke model", "error", err.Error(), "model", *withModel, "data", withData)
		return nil, err
	}

	var res B
	err = json.Unmarshal(out.Body, &res)
	if err != nil {
		slog.Debug("failed to unmarshal response", "error", err.Error(), "model", *withModel, "data", withData)
		return nil, err
	}

	return &Response{
		Completion: withParser(res),
	}, nil
}
