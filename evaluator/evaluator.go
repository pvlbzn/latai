package evaluator

import (
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/pvlbzn/latai/prompt"
	"github.com/pvlbzn/latai/provider"
)

var (
	ErrNoProvider = errors.New("no provider specified")
	ErrNoModel    = errors.New("no model provided")
	ErrNoPrompt   = errors.New("no prompt(s) provided")
	ErrSampleSize = errors.New("sample size must be 1 or more")
)

type Evaluator struct {
	provider   provider.Provider
	model      *provider.Model
	prompts    []*prompt.Prompt
	sampleSize int
	// concurrency
	// timeout
}

type Evaluation struct {
	ModelName     string
	ModelProvider string
	Responses     []string
	Latency       time.Duration
}

func NewEvaluator(provider provider.Provider, model *provider.Model, prompts ...*prompt.Prompt) *Evaluator {
	return &Evaluator{
		provider:   provider,
		model:      model,
		sampleSize: len(prompts),
		prompts:    prompts,
	}
}

func (e *Evaluator) WithSampleSize(n int) *Evaluator {
	e.sampleSize = n
	return e
}

func (e *Evaluator) validate() error {
	if e.provider == nil {
		return ErrNoProvider
	}

	if e.model == nil {
		return ErrNoModel
	}

	if e.sampleSize <= 0 {
		return ErrSampleSize
	}

	if e.prompts == nil || len(e.prompts) == 0 {
		return ErrNoPrompt
	}

	return nil
}

// Evaluate a model with provided prompts of a given sample size. Provided model
// becomes model under test. Model will be tested against prompts using provided
// number of samples.
//
// There is a direct relation between prompts and sample size.  When count of provided
// prompts equals to sample size then each prompt runs once. This is the recommended
// and default behavior. In this way measurement stays more  precise because one sample
// is not enough to define latency, however, one prompt shouldn't be re-used for clean
// measurement due to possible prompt caching internal mechanisms.
//
// If required, sample size may be changed using `Evaluator.WithSampleSize`. Evaluate
// will detect that amount of prompts doesn't match sample size and will run sampling
// picking up a random prompt from the prompt pool. This measurement might be affected
// by prompt caching.
func (e *Evaluator) Evaluate() (*Evaluation, error) {
	// Validate.
	err := e.validate()
	if err != nil {
		slog.Debug("failed to run evaluator", "error", err.Error())
		return nil, err
	}

	// Get metrics.
	var metrics []*provider.Metric

	if len(e.prompts) != e.sampleSize {
		metrics, err = e.runRandomSample()
	} else {
		metrics, err = e.runUniqueSample()
	}
	if err != nil {
		slog.Debug("failed to run a sample", "error", err.Error())
		return nil, err
	}

	// Combine.
	var responses []string
	var latency time.Duration
	for _, m := range metrics {
		latency += m.Latency
		responses = append(responses, m.Response.Completion)
	}
	latency /= time.Duration(len(metrics))

	return &Evaluation{
		ModelName:     e.model.Name,
		ModelProvider: string(e.model.Provider),
		Responses:     responses,
		Latency:       latency,
	}, nil
}

// runUniqueSample runs measurements which are unique and may defeat prompt caching.
func (e *Evaluator) runUniqueSample() ([]*provider.Metric, error) {
	var res []*provider.Metric

	for _, p := range e.prompts {
		m, err := e.provider.Measure(e.model, p)
		if err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	return res, nil
}

// runRandomSample runs measurements picking up prompts randomly out of prompt pool.
func (e *Evaluator) runRandomSample() ([]*provider.Metric, error) {
	var res []*provider.Metric

	for i := 0; i < e.sampleSize; i++ {
		randomPrompt := e.prompts[rand.Intn(len(e.prompts))]
		m, err := e.provider.Measure(e.model, randomPrompt)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}

	return res, nil
}
