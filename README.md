# `latai` Generative AI Latency Measurement TUI

# Usage

Latai currently supports following providers:
* OpenAI
* Groq
* Bedrock

If you need to add other providers and/or models read Contribution section.

## API Keys
To access LLMs Latai has to have access to API keys. Each provider is optional. By default, Latai loads all providers and verifies their keys. If keys aren't found provider is not loaded. Therefore, if you don't need some provider just don't add its key. 

TLDR: add following keys and values into your environment and update your terminal environment.

```shell
# OpenAI API key.
export OPENAI_API_KEY=

# Groq API key.
export GROQ_API_KEY=

# AWS Bedrock key. You can specify your AWS profile and region
# here. If you don't do this, yet you have your AWS CLI installed
# Latai will use `default` profile and `us-east-1` region.
export AWS_PROFILE=
export AWS_REGION=
```

> [!IMPORTANT]
> Transparency note. Keys never leave your machine. Latai has no telemetry and does not send your data anywhere. Each provider code has two functions which reach internet. First one is `VerifyAccess` which sends request to list all available models to check API key validity. Second one is `Send` which sends requests to LLMs to measure latency based on default or user prompts.

API key management is provider specific, here are the instructions for each supported provider.

### OpenAI, Groq

OpenAI and Groq use the same API therefore their key management principle is identical. To set keys add these into your environment:

```shell
# OpenAI API key.
export OPENAI_API_KEY=

# Groq API key.
export GROQ_API_KEY=
```

If you don't need Groq, just don't add key.


### AWS Bedrock

AWS uses their own mechanism of authentication which is based on [AWS CLI](https://aws.amazon.com/cli/). Refer to their documentation for details if you need it.

Latai will load your AWS profile in following order:
1. Access `AWS_PROFILE` and `AWS_REGION` from your environment.
2. If not found default to the default values `AWS_PROFILE=default`, `AWS_REGION=us-east-1`.

To set your profile and region add those:

```shell
export AWS_PROFILE=
export AWS_REGION=
```

Make sure that you either load Latai from the same terminal after exports, or add those exports into your shell `rc` file, e.g. `.bashrc`, `.zshrc`, etc.

> [!NOTE]
> Make sure you have access to LLM models from your AWS account. They are not enabled by default. You have to navigate to `https://REGION.console.aws.amazon.com/bedrock/home?region=REGION#/modelaccess` and enable models from the console. Make sure to replace `REGION` with your actual region. [Here is the link](https://us-east-1.console.aws.amazon.com/bedrock/home?region=us-east-1#/modelaccess) for `us-east-1`.

To verify your access you can use `aws` CLI.

```shell
aws bedrock \
  list-foundation-models \
  --region REGION \
  --profile PROFILE
```

Substitute `REGION` and `PROFILE` with your data. You can optionally pipe into `jq` to make output more readable.


## Prompts: Default and Custom

Latai uses a set 3 pre-defined prompts by default. They are just good enough to measure latency to model and back. E.g. `Respond with a single word: "optimistic".`. You can find them [here](https://github.com/pvlbzn/latai/tree/main/internal/prompt/prompts). Three pre-defined prompts meaning that by default all sampling happens with 3 runs.

If you  need to measure compute time, or performance with your particular prompts, then you can add your own into `~/.latai` directory.

```shell
# Create a directory where prompts are stored.
mkdir -p ~/.latai/prompts

# Create your prompts in there.
cd ~/.latai/prompts
touch p1.prompt

# Or create many.
touch {p1, p2, p3, p4, p5}.prompt
```

You can create any number of prompts you wish, just mind throttling and rate limiting. All prompts should have `.prompt` postfix, files with other postfixes will be ignored.



# Latency Measurement Strategy

**TODO: add description of work**

**TODO: add prompts and count their token load for measurements**

LLM providers tend to have a significant latency jitter. Latency varies greatly between calls, and occasionally request may timeout. To provide  a stable means of measurement _multi-sample_ approach is used.

Another point of consideration is prompt caching. Major providers started to roll out this feature around Q3 2024.

**TODO: IMPLEMENT**\
To make latency measurements as close to real-world as possible Genlat provides a collection of prompts.

> [!NOTE]
> Prompt caching usually last up to 15 minutes. This is not set in stone number, therefore take this into consideration if running measurements repeatedly.

Collection of prompts can be found at `./prompt/prompts` directory. They are designed to be similar to generate comparable load on LLM, yet different enough to not get cached.

By default, Genlat uses sampling size of 3 for the measurements. That means that it will pick 3 prompts from prompts directory and run them against each LLM. 

Prompt caching:
- [OpenAI prompt caching](https://platform.openai.com/docs/guides/prompt-caching) [^1]
- [Anthropic prompt caching](https://docs.anthropic.com/en/docs/build-with-claude/prompt-caching) [^2]

[^1]: OpenAI claims up to 80% reduced latency with prompt caching.
[^2]: OpenAI claims up to 85% reduced latency with prompt caching for long prompts.


## Customizing Prompts

**TODO: IMPLEMENT**\
You can provide your own prompts by creating `~/.genlat/prompts` directory and putting prompts there. This can be useful for custom use cases where prompts or / and use cases differ from standard ones. Doing so be wary of prompt caching.


# Providers & Vendors & Models

Definitions:
* Provider is a service provider which serves access to model(s) over some API.
* Vendor is a company-owner of a model which created, trained, and aligned a model.
* Model is a LLM model with particular properties such as performance, context, languages, etc.


Providers are services which serve models over API. Models can be separated by families, for example Claude family of models. Some providers are mono-family, e.g. OpenAI, and follow single unified API for all underlying models. Other providers are multi-family, e.g. AWS Bedrock. Multi-family providers have their API, however particular format of communication depends on the model's family. Vendors may have one or more families of models. Families generally defined by their API, for example if model A and B have the same API and belong to the same vendor then they belong to the same family.

Genlat aggregates models by provider, because provider is the root of the relation, and it is what runs models. However, more frequently than not a model can run on a multiple providers. Thus, there are two APIs - provider API, and model API. To simplify provider API is responsible for transport layer, and model API is responsible for data format. 


## Rate Limits
Commonly rate limits measured in following metrics:
* RPM: Requests per minute
* RPD: Requests per day
* TPM: Tokens per minute
* TPD: Tokens per day

Verify those with your model provider. This information can be found at provider's documentation. You can find these links below. Keep in mind that these rate limits almost always negotiable with your provider, and generally limits applied to models, not provider itself.

Providers:
* [Groq Rate Limits](https://console.groq.com/docs/rate-limits)
* [OpenAI Rate Limits](https://platform.openai.com/docs/guides/rate-limits)
* AWS Bedrock (read below)


### OpenAI

OpenAI uses tiered rate limits from 1 to 5. For more details consult their documentation.

### AWS Bedrock

AWS Bedrock has multiple providers under their name. Before using most of the models you have to request access to them via AWS UI. 

#### Access

> [!NOTE]
> Make sure you have access to LLM models from your AWS account. They are not enabled by default. You have to navigate to https://REGION.console.aws.amazon.com/bedrock/home?region=REGION#/modelaccess and enable models from the console. Make sure to replace `REGION` with your actual region.

To verify your access you can use `aws` CLI.

```shell
aws bedrock \
  list-foundation-models \
  --region REGION \
  --profile PROFILE
```

Substitute `REGION` and `PROFILE` with your data. You can optionally pipe into `jq` to make output more readable.


#### Models

Even though AWS Bedrock returns lots of models, not all of them can be accessed "as-is". For example AWS Bedrock lists more than 20 Claude-family models, however, only 6 out of them are available without provisioning. Genlat doesn't include models which require special access at the moment.

You can fork this repository and add required provisioned models by adding their ID into `NewBedrock` function [in this file](internal/provider/bedrock.go).


# Contributing

## Adding a New Provider

All providers reside at [provider package](internal/provider). There is one main interface defined at [`provider.go`](internal/provider/provider.go) file.

```go
// Comments omitted for brevity, check source file 
// to see full version.
type Provider interface {
	Name() ModelProvider
	GetLLMModels(filter string) []*Model
	Measure(model *Model, prompt *prompt.Prompt) (*Metric, error)
	Send(message string, to *Model) (*Response, error)
	VerifyAccess() bool
}
```

If a struct satisfies this `Provider` interface it is ready to be used along with all other providers.

If provider you are adding is OpenAI API compliant check [`groq.go`](internal/provider/groq.go) implementation.

Do not forget to add tests. You can see implementation of tests inside of [`provider` package](internal/provider).



# Troubleshooting

## `err` as Latency Value

Read `Events` block of TUI, it generally explains what went wrong. The most common issues is related to AWS Bedrock due to access to models.
