package main

import (
	"fmt"
	"log/slog"

	"github.com/pvlbzn/genlat/evaluator"
	"github.com/pvlbzn/genlat/prompt"
	"github.com/pvlbzn/genlat/provider"
	"github.com/pvlbzn/genlat/tui"
)

func main() {
	tui.Run()
	//setLogger(slog.LevelError)
	//
	//if err := runOpenAI(); err != nil {
	//	panic(err)
	//}
}

func setLogger(level slog.Level) {
	slog.SetLogLoggerLevel(level)
}

func runOpenAI() error {
	p, err := provider.NewOpenAI("")
	if err != nil {
		return err
	}

	models, err := p.GetLLMModels("")
	if err != nil {
		return err
	}

	for _, m := range models {
		_, err := sendMessage(p, m)
		if err != nil {
			slog.Error("message error", "error", err, "id", m.ID)
		}

		prompts, err := prompt.GetPrompts()
		if err != nil {
			return err
		}

		eval := evaluator.NewEvaluator(p, m, prompts...)

		metrics, err := eval.Evaluate()
		if err != nil {
			return err
		}

		fmt.Printf("%+v", metrics)
	}

	return nil
}

//func runBedrock() error {
//	p, err := provider.NewBedrock(provider.DefaultBedrockRegion)
//	if err != nil {
//		return err
//	}
//
//	models, err := listModels(p, "")
//	if err != nil {
//		return err
//	}
//
//	if err = sendMessage(p, models[2]); err != nil {
//		return err
//	}
//
//	return nil
//}

//func renderUI() {
//	p := tea.NewProgram(tui.InitialModel())
//	if _, err := p.Run(); err != nil {
//		fmt.Println("Error:", err)
//		os.Exit(1)
//	}
//}

func listModels(b provider.Provider, filter string) ([]*provider.Model, error) {
	models, err := b.GetLLMModels(filter)
	if err != nil {
		return nil, err
	}

	for _, model := range models {
		fmt.Printf("model: %s, id: %s, vendor: %s, provider: %s\n", model.Name, model.ID, model.Vendor, model.Provider)
	}

	return models, nil
}

func sendMessage(b provider.Provider, m *provider.Model) (string, error) {
	res, err := b.Send("What is your name? Reply in a single word.", m)
	if err != nil {
		return "", err
	}

	return res, nil
}
