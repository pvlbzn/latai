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

type Region = string

const (
	DefaultBedrockRegion Region = "us-east-1"
)

type BedrockModel struct {
	ID       string
	Provider string
	Name     string
}

type Bedrock struct {
	client  *bedrock.Client
	runtime *bedrockruntime.Client
}

func NewBedrock(region Region) (*Bedrock, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		slog.Error("failed to load bedrock default configuration", "error", err)
		return nil, err
	}

	return &Bedrock{
		client:  bedrock.NewFromConfig(cfg),
		runtime: bedrockruntime.NewFromConfig(cfg),
	}, nil
}

type claudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type claudeResponse struct {
	Completion string `json:"completion"`
}

const claudePromptFormat = "\n\nHuman: %s\n\nAssistant:"

// GetModels from AWS Bedrock service. Effectively lists available models.
func (s *Bedrock) GetModels(filter string) ([]*Model, error) {
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
func (s *Bedrock) Send(message string, to *Model) (string, error) {
	slog.Debug("sending message", "message", message, "model", *to)

	// TODO: not all models work with documented API, therefore for now
	// 	this value is hardcoded.
	modelID := "anthropic.claude-v2"

	data := claudeRequest{
		Prompt:            fmt.Sprintf(claudePromptFormat, message),
		MaxTokensToSample: 2048,
		Temperature:       1.0,
		StopSequences:     []string{"\n\nHuman:"},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal model data", "error", err.Error(), "data", data)
		return "", err
	}

	out, err := s.runtime.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Body:        dataBytes,
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		slog.Error("failed to invoke model", "error", err.Error(), "model", *to, "data", data)
		return "", err
	}

	var res claudeResponse
	err = json.Unmarshal(out.Body, &res)
	if err != nil {
		slog.Error("failed to unmarshal response", "error", err.Error(), "model", *to, "data", data)
		return "", err
	}

	return res.Completion, nil
}
