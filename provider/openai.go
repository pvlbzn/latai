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

	return &OpenAI{client: c}, nil
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

func (s *OpenAI) GetModels(filter string) ([]*Model, error) {
	res, err := s.client.ListModels(context.TODO())
	if err != nil {
		slog.Error("failed to get list of models", "error", err)
		return nil, err
	}

	models := make([]*Model, 0, len(res.Models))

	for _, model := range res.Models {
		modelName, query := strings.ToLower(model.ID), strings.ToLower(filter)

		if strings.Contains(modelName, query) {
			models = append(models, &Model{
				ID:       model.ID,
				Name:     s.id2Name(model.ID),
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

func (s *OpenAI) Send(message string, to *Model) (string, error) {
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
		return "", err
	}

	return res.Choices[0].Message.Content, nil
}
