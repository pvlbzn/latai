package prompt

import (
	"embed"
	"log/slog"
	"path/filepath"
	"strings"
)

// Prompt structure which represents a single prompt.
type Prompt struct {
	// Description is a free form description of the prompt.
	Description string

	// Content of the prompt.
	Content string
}

//go:embed prompts/*.prompt
var defaultPrompts embed.FS

// GetPrompts returns prompts for evaluation. Returns either user-defined prompts
// from `~/.latai/prompts/*.prompt`, or default embedded prompts.
// TODO: user defined prompts.
func GetPrompts() ([]*Prompt, error) {
	// TODO: first load files from ~/.latai/prompt/*, then if empty
	// 	load default prompts.
	return loadDefaultPrompts()
}

func loadDefaultPrompts() ([]*Prompt, error) {
	var prompts []*Prompt

	files, err := defaultPrompts.ReadDir("prompts")
	if err != nil {
		slog.Error("failed to read defaultPrompts", "err", err)
		return nil, err
	}

	for _, f := range files {
		// Ignore everything but `*.prompt` files.
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".prompt") {
			continue
		}

		content, err := defaultPrompts.ReadFile(filepath.Join("prompts", f.Name()))
		if err != nil {
			slog.Error("failed get content of a file", "err", err, "file", f.Name())
			return nil, err
		}

		prompts = append(prompts, &Prompt{
			Description: "Default prompt " + f.Name(),
			Content:     string(content),
		})
	}

	return prompts, nil
}
