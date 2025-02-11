package prompt

import (
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type PromptType string

const (
	PromptTypeUser    PromptType = "user"
	PromptTypeDefault PromptType = "default"
)

// Prompt structure which represents a single prompt.
type Prompt struct {
	// Type of prompt, either default or user.
	Type PromptType

	// Description is a free form description of the prompt.
	Description string

	// Content of the prompt.
	Content string
}

//go:embed prompts/*.prompt
var defaultPrompts embed.FS

// GetPrompts returns prompts for evaluation. Returns either user-defined prompts
// from `~/.latai/prompts/*.prompt`, or default embedded prompts.
func GetPrompts() ([]*Prompt, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".latai", "prompts")

	prompts, err := loadUserPrompts(dir)
	if err != nil || len(prompts) == 0 {
		return loadDefaultPrompts()
	}

	return prompts, nil
}

// loadUserPrompts loads prompt files from a given directory.
func loadUserPrompts(dir string) ([]*Prompt, error) {
	var prompts []*Prompt

	files, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".prompt" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			continue
		}

		prompts = append(prompts, &Prompt{
			Type:        PromptTypeUser,
			Description: "User prompt " + file.Name(),
			Content:     string(content),
		})
	}

	return prompts, nil
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
			Type:        PromptTypeDefault,
			Description: "Default prompt " + f.Name(),
			Content:     string(content),
		})
	}

	return prompts, nil
}
