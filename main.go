package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pvlbzn/genlat/evaluator"
	"github.com/pvlbzn/genlat/prompt"
	"github.com/pvlbzn/genlat/provider"
	"github.com/pvlbzn/genlat/tui"
)

func main() {
	setLogger()

	if err := runOpenAI(); err != nil {
		panic(err)
	}
}

func setLogger() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func runOpenAI() error {
	p, err := provider.NewOpenAI("")
	if err != nil {
		return err
	}

	models, err := listModels(p, "gpt-4o-mini-2024-07-18")
	if err != nil {
		return err
	}

	model := models[0]

	if err = sendMessage(p, model); err != nil {
		return err
	}

	prompts, err := prompt.GetPrompts()
	if err != nil {
		return err
	}

	for _, prompt := range prompts {
		fmt.Printf("Prompt: %+v\n", prompt)
	}

	eval := evaluator.NewEvaluator(p, model, prompts...)

	metrics, err := eval.Evaluate()
	if err != nil {
		return err
	}

	fmt.Printf("%+v", metrics)

	return nil
}

func runBedrock() error {
	p, err := provider.NewBedrock(provider.DefaultBedrockRegion)
	if err != nil {
		return err
	}

	models, err := listModels(p, "claude")
	if err != nil {
		return err
	}

	if err = sendMessage(p, models[2]); err != nil {
		return err
	}

	return nil
}

func renderUI() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func listModels(b provider.Provider, filter string) ([]*provider.Model, error) {
	models, err := b.GetModels(filter)
	if err != nil {
		return nil, err
	}

	for _, model := range models {
		fmt.Printf("model: %s, id: %s, vendor: %s, provider: %s\n", model.Name, model.ID, model.Vendor, model.Provider)
	}

	return models, nil
}

func sendMessage(b provider.Provider, m *provider.Model) error {
	res, err := b.Send("What is your name? Reply in a single word.", m)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
