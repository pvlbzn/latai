package main

import (
	"fmt"

	"github.com/pvlbzn/genlat/provider"
)

func main() {
	//p := tea.NewProgram(ui.InitialModel())
	//if _, err := p.Run(); err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(1)
	//}

	b := provider.NewBedrock()
	models, err := b.GetModelsList()
	if err != nil {
		panic(err)
	}

	for _, model := range models {
		fmt.Printf("model: %s, provider: %s\n", model.Name, model.Provider)
	}
}
