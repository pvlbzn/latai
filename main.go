package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pvlbzn/genlat/provider"
	"github.com/pvlbzn/genlat/ui"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	b, err := provider.NewBedrock(provider.DefaultBedrockRegion)
	if err != nil {
		return err
	}

	if err = listModels(b); err != nil {
		return err
	}

	if err = sendMessage(b); err != nil {
		return err
	}

	return nil
}

func renderUI() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func listModels(b *provider.Bedrock) error {
	models, err := b.GetModels()
	if err != nil {
		return err
	}

	for _, model := range models {
		fmt.Printf("model: %s, id: %s, provider: %s\n", model.Name, model.ID, model.Provider)
	}

	return nil
}

func sendMessage(b *provider.Bedrock) error {
	res, err := b.Send("What is your name?")
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
