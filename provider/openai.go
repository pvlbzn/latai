package provider

import (
	"context"
	"errors"
	"github.com/pvlbzn/latai/prompt"
	"github.com/sashabaranov/go-openai"
	"log/slog"
	"os"
	"strings"
)

type OpenAI struct {
	client *openai.Client
	models []Model
}

func NewOpenAI(apiKey string) (*OpenAI, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, errors.New("openai api key not found")
		}
	}

	c := openai.NewClient(apiKey)
	models := []Model{
		{ID: "gpt-4-1106-preview", Name: "GPT 4 1106 Preview", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-3.5-turbo", Name: "GPT 3.5 Turbo", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-3.5-turbo-0125", Name: "GPT 3.5 Turbo 0125", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "o1-mini", Name: "O1 Mini", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "o1-mini-2024-09-12", Name: "O1 Mini 2024 0`9 12", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "o1-2024-12-17", Name: "O1 2024 12 17", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-3.5-turbo-16k", Name: "GPT 3.5 Turbo 16k", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "o1", Name: "O1", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "o1-preview-2024-09-12", Name: "O1 Preview 2024 09 12", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "o1-preview", Name: "O1 Preview", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4", Name: "GPT 4", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4-0613", Name: "GPT 4 0613", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "chatgpt-4o-latest", Name: "ChatGPT 4o Latest", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4o-2024-08-06", Name: "GPT 4o 2024 08 06", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4o", Name: "GPT 4o", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-3.5-turbo-1106", Name: "GPT 3.5 Turbo 1106", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4-turbo-2024-04-09", Name: "GPT 4 Turbo 2024 04 09", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4-turbo", Name: "GPT 4 Turbo", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4-turbo-preview", Name: "GPT 4 Turbo Preview", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4o-2024-05-13", Name: "GPT 4o 2024 05 13", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4o-2024-11-20", Name: "GPT 4o 2024 11 20", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4o-mini-2024-07-18", Name: "GPT 4o Mini 2024 07 18", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4o-mini", Name: "GPT 4o Mini", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
		{ID: "gpt-4-0125-preview", Name: "GPT 4 0125 Preview", Provider: ModelProviderOpenAI, Vendor: ModelVendorOpenAI, Family: ModelFamilyGPT},
	}

	return &OpenAI{client: c, models: models}, nil
}

// id2Name converts ID to name format. E.g. transforms `o1-preview-2024-09-12`
// to `o1 preview 2024 09 12`.
func (s *OpenAI) id2Name(id string) string {
	elems := strings.Split(id, "-")
	for i := range elems {
		elems[i] = strings.ToUpper(elems[i][:1]) + elems[i][1:]
	}
	return strings.Join(elems, " ")
}

// GetLLMModels returns LLM models only. Filter is applied to search
// models by their name, e.g. "4o" filter will return 4o family models.
// Empty filter returns full list of available LLM models.
func (s *OpenAI) GetLLMModels(filter string) []*Model {
	return filterModels(s.models, filter)
}

func (s *OpenAI) Send(message string, to *Model) (*Response, error) {
	slog.Debug("sending message", "message", message, "to", to)

	res, err := s.client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: to.ID,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: message},
			},
		})

	if err != nil {
		return nil, err
	}

	return &Response{Completion: res.Choices[0].Message.Content}, nil
}

func (s *OpenAI) Measure(model *Model, prompt *prompt.Prompt) (*Metric, error) {
	return measure(s, model, prompt)
}
