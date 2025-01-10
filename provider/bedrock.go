package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
)

type Bedrock struct {
}

func NewBedrock() *Bedrock {
	return &Bedrock{}
}

type BedrockModel struct {
	Provider string
	Name     string
}

func (b *Bedrock) GetModelsList() ([]*BedrockModel, error) {
	region := "us-east-1"

	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		slog.Error("couldn't load default configuration", "error", err.Error())
		return nil, err
	}

	client := bedrock.NewFromConfig(sdkConfig)
	res, err := client.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
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
			&BedrockModel{Provider: *summary.ProviderName, Name: *summary.ModelName})
	}

	return models, nil
}
