package provider

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/pvlbzn/genlat/prompt"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client *openai.Client
	models []Model
}

func NewOpenAI(apiKey string) (*OpenAI, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			m := "openai api key not found"
			slog.Error(m)
			return nil, errors.New(m)
		}
	}

	c := openai.NewClient(apiKey)
	models := []Model{
		{ID: "gpt-4-1106-preview", Name: "GPT 4 1106 Preview", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-3.5-turbo", Name: "GPT 3.5 Turbo", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-3.5-turbo-0125", Name: "GPT 3.5 Turbo 0125", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "o1-mini", Name: "O1 Mini", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "o1-mini-2024-09-12", Name: "O1 Mini 2024 0`9 12", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "o1-2024-12-17", Name: "O1 2024 12 17", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-3.5-turbo-16k", Name: "GPT 3.5 Turbo 16k", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "o1", Name: "O1", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "o1-preview-2024-09-12", Name: "O1 Preview 2024 09 12", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "o1-preview", Name: "O1 Preview", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4", Name: "GPT 4", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4-0613", Name: "GPT 4 0613", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "chatgpt-4o-latest", Name: "ChatGPT 4o Latest", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4o-2024-08-06", Name: "GPT 4o 2024 08 06", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4o", Name: "GPT 4o", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-3.5-turbo-1106", Name: "GPT 3.5 Turbo 1106", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4-turbo-2024-04-09", Name: "GPT 4 Turbo 2024 04 09", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4-turbo", Name: "GPT 4 Turbo", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4-turbo-preview", Name: "GPT 4 Turbo Preview", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4o-2024-05-13", Name: "GPT 4o 2024 05 13", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4o-2024-11-20", Name: "GPT 4o 2024 11 20", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4o-mini-2024-07-18", Name: "GPT 4o Mini 2024 07 18", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4o-mini", Name: "GPT 4o Mini", Provider: "OpenAI", Vendor: "OpenAI"},
		{ID: "gpt-4-0125-preview", Name: "GPT 4 0125 Preview", Provider: "OpenAI", Vendor: "OpenAI"},
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
func (s *OpenAI) GetLLMModels(filter string) ([]*Model, error) {
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

// GetAllModels returns models fetched from OpenAI API. Note that
// these models are mixed, some models are embeddings, some LLM, some
// multimodal.
func (s *OpenAI) GetAllModels(filter string) ([]*Model, error) {
	res, err := s.client.ListModels(context.TODO())
	if err != nil {
		slog.Error("failed to get list of models", "error", err)
		return nil, err
	}

	models := make([]*Model, 0, len(res.Models))

	for _, model := range res.Models {
		modelName, query := s.id2Name(strings.ToLower(model.ID)), strings.ToLower(filter)

		if strings.Contains(modelName, query) {
			models = append(models, &Model{
				ID:       model.ID,
				Name:     modelName,
				Provider: "OpenAI",
				Vendor:   "OpenAI",
			})
		}
	}

	return models, nil
}

func (s *OpenAI) Measure(model *Model, prompt *prompt.Prompt) (*Metric, error) {
	start := time.Now()
	res, err := s.Send(prompt.Content, model)
	elapsed := time.Since(start)
	if err != nil {
		return nil, err
	}

	return &Metric{
		Model:    model,
		Latency:  elapsed,
		Response: res,
	}, nil
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
		slog.Error("failed to send openai message", "error", err, "message", message)
		return nil, err
	}

	return &Response{Completion: res.Choices[0].Message.Content}, nil
}
