package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/pvlbzn/latai/prompt"
	"github.com/sashabaranov/go-openai"
	"os"
)

type Groq struct {
	client *openai.Client
	models []Model
}

// NewGroq initializes and returns a new Groq instance.
func NewGroq(apiKey string) (*Groq, error) {
	if apiKey == "" {
		apiKey = os.Getenv("GROQ_API_KEY")
		if apiKey == "" {
			return nil, errors.New("groq api key not found")
		}
	}

	conf := openai.DefaultConfig(apiKey)
	conf.BaseURL = "https://api.groq.com/openai/v1"
	c := openai.NewClientWithConfig(conf)

	models := []Model{
		{ID: "gemma2-9b-it", Name: "Gemma 2 9B IT", Provider: ModelProviderGroq, Vendor: ModelVendorGoogle, Family: ModelFamilyGemma},
		{ID: "llama-3.3-70b-versatile", Name: "Llama 3.3 70b Versatile", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
		{ID: "llama-3.1-8b-instant", Name: "Llama 3.1 8b Instant", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
		{ID: "llama-guard-3-8b", Name: "Llama Guard 3 8B", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
		{ID: "llama3-70b-8192", Name: "Llama3 70b 8192", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
		{ID: "llama3-8b-8192", Name: "Llama3 8b 8192", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
		{ID: "mixtral-8x7b-32768", Name: "Mixtral 8x7b 32768", Provider: ModelProviderGroq, Vendor: ModelVendorMistralAI, Family: ModelFamilyMixtral},
		{ID: "deepseek-r1-distill-llama-70b", Name: "DeepSeek R1 Distill Llama 70B", Provider: ModelProviderGroq, Vendor: ModelVendorDeepSeek, Family: ModelFamilyR1},
		{ID: "llama-3.2-1b-preview", Name: "Llama 3.2 1b Preview", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
		{ID: "llama-3.2-3b-preview", Name: "Llama 3.2 3b Preview", Provider: ModelProviderGroq, Vendor: ModelVendorMeta, Family: ModelFamilyLlama3},
	}

	return &Groq{
		client: c,
		models: models,
	}, nil
}

func (s *Groq) GetLLMModels(filter string) []*Model {
	return filterModels(s.models, filter)
}

func (s *Groq) Send(message string, model *Model) (*Response, error) {
	switch model.Vendor {
	case ModelVendorGoogle:
		return s.runGroqInference(model, message)
	case ModelVendorMeta:
		return s.runGroqInference(model, message)
	case ModelVendorMistralAI:
		return s.runGroqInference(model, message)
	case ModelVendorDeepSeek:
		return s.runGroqInference(model, message)

	default:
		return nil, fmt.Errorf("unsupported vendor: %s", model.Vendor)
	}
}

func (s *Groq) runGroqInference(model *Model, message string) (*Response, error) {
	res, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model.ID,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: message},
			},
		})
	if err != nil {
		return nil, err
	}

	return &Response{Completion: res.Choices[0].Message.Content}, nil
}

//func (s *Groq) runGroqInferenceGemmaFamily(message string, model *Model) (*Response, error) {
//
//}
//
//func (s *Groq) runGroqInferenceLlamaFamily(message string, model *Model) (*Response, error) {
//	res, err := s.client.CreateChatCompletion(
//		context.Background(),
//		openai.ChatCompletionRequest{
//			Model: model.ID,
//			Messages: []openai.ChatCompletionMessage{
//				{Role: openai.ChatMessageRoleUser, Content: message},
//			},
//		})
//	if err != nil {
//		return nil, err
//	}
//
//	return &Response{Completion: res.Choices[0].Message.Content}, nil
//}
//
//func (s *Groq) runGroqInferenceMixtralFamily(message string, model *Model) (*Response, error) {
//
//}
//
//func (s *Groq) runGroqInferenceR1Family(message string, model *Model) (*Response, error) {
//
//}
//
//func runGroqInferece[A, B any](groq *Groq, withModel *Model, withData A, withParser func(B) string) (*Response, error) {
//	res, err := groq.client
//}

func (s *Groq) Measure(model *Model, prompt *prompt.Prompt) (*Metric, error) {
	return measure(s, model, prompt)
}
