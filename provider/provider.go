package provider

type Prompt struct {
	Description string
	Content     string
}

type Metrics struct {
	Provider string
	Model    string
	Latency  string
	Tokens   string
}

type Provider interface {
	Run() error
}
