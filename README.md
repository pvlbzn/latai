# Generative AI Latency Measurement


## Latency Measurement Strategy

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


### Customizing Prompts

**TODO: IMPLEMENT**\
You can provide your own prompts by creating `~/.genlat/prompts` directory and putting prompts there. This can be useful for custom use cases where prompts or / and use cases differ from standard ones. Doing so be wary of prompt caching.


## Providers & Models

### OpenAI
gpt4o
gpt4o-mini

### AWS Bedrock
Anthropic

### Microsoft
https://huggingface.co/microsoft/phi-4
ollama
