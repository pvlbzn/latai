package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Region = string

const (
	DefaultBedrockRegion Region = "us-east-1"
)

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

type BedrockModel struct {
	ID       string
	Provider string
	Name     string
}

type request struct {
	Prompt            string `json:"prompt"`
	MaxTokensToSample int    `json:"max_tokens_to_sample"`
}

type response struct {
	Completion string `json:"completion"`
}

const claudePromptFormat = "\n\nHuman: %s\n\nAssistant:"

// GetModels from AWS Bedrock service. Effectively lists available models.
func (b *Bedrock) GetModels() ([]*BedrockModel, error) {
	res, err := b.client.ListFoundationModels(context.TODO(), &bedrock.ListFoundationModelsInput{})
	if err != nil {
		slog.Error("couldn't get list of foundation models", "error", err.Error())
		return nil, err
	}

	if len(res.ModelSummaries) == 0 {
		m := "no foundation models found"
		slog.Error(m)
		return nil, fmt.Errorf(m)
	}

	models := make([]*BedrockModel, 0, len(res.ModelSummaries))
	for _, summary := range res.ModelSummaries {
		models = append(
			models,
			&BedrockModel{
				ID:       *summary.ModelId,
				Provider: *summary.ProviderName,
				Name:     *summary.ModelName})
	}

	return models, nil
}

// Send message.
func (b *Bedrock) Send(message string) (string, error) {
	// TODO: add model / provider / id
	model := "anthropic.claude-v2"

	slog.Debug("sending message", "message", message, "")

	data := request{Prompt: fmt.Sprintf(claudePromptFormat, message), MaxTokensToSample: 2048}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal model data", "error", err.Error(), "data", data)
		return "", err
	}

	out, err := b.runtime.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Body:        dataBytes,
		ModelId:     aws.String(model),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		slog.Error("failed to invoke model", "error", err.Error(), "model", model, "data", data)
		return "", err
	}

	var res response
	err = json.Unmarshal(out.Body, &res)
	if err != nil {
		slog.Error("failed to unmarshal response", "error", err.Error(), "model", model, "data", data)
		return "", err
	}

	return res.Completion, nil
}
