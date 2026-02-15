# exa-cli

CLI for the [Exa AI](https://exa.ai) search API.

Search the web, find similar pages, get AI-powered answers with citations, retrieve page contents, and explore code context — all from the command line.

## Installation

```bash
brew install roboalchemist/tap/exa-cli
```

Or build from source:

```bash
git clone https://github.com/roboalchemist/exa-cli.git
cd exa-cli
make build
sudo make install
```

## Authentication

Get your API key at [dashboard.exa.ai](https://dashboard.exa.ai/api-keys).

```bash
# Option 1: Environment variable (recommended)
export EXA_API_KEY="your-api-key"

# Option 2: Config file
exa auth
```

## Usage

### Search

```bash
# Basic search
exa search "hottest AI startups"

# Deep search with more results
exa search "climate change research" --type deep -n 20

# Filter by domain and date
exa search "machine learning" --include-domains arxiv.org --start-date 2025-01-01

# Search with highlights and summary
exa search "quantum computing" --highlights --summary

# Category-specific search
exa search "OpenAI" --category company
```

### Find Similar

```bash
# Find pages similar to a URL
exa similar "https://arxiv.org/abs/2307.06435" -n 10

# Exclude the source domain
exa similar "https://blog.example.com" --exclude-source
```

### AI Answers

```bash
# Get an answer with citations
exa answer "What is the capital of France?"

# Stream the answer
exa answer "Explain quantum computing" --stream

# Include full source text
exa answer "Latest AI breakthroughs" --text
```

### Page Contents

```bash
# Get page text
exa contents https://example.com

# With highlights and summary
exa contents https://example.com --highlights --summary

# Multiple URLs
exa contents https://a.com https://b.com
```

### Code Context

```bash
# Get code-relevant context
exa context "React hooks state management"

# Limit tokens
exa context "Go error handling" --tokens 5000
```

### Usage & Billing

```bash
exa usage
exa usage --start-date 2025-01-01
```

## Output Formats

All commands support multiple output formats:

```bash
# Table (default) — human-readable
exa search "AI" -n 3

# JSON — structured, agent-friendly
exa search "AI" -n 3 --json

# JSON with field selection
exa search "AI" --json --fields title,url,score

# JSON with jq filtering
exa search "AI" --json --jq '.results[] | {title, url}'

# Plaintext — tab-separated for piping
exa search "AI" -n 3 --plaintext
```

## Search Types

| Type | Description | Latency |
|------|-------------|---------|
| `auto` | Combines methods with reranker (default) | Medium |
| `fast` | Optimized for speed | <400ms |
| `deep` | Query expansion, comprehensive | Higher |
| `neural` | Pure embeddings-based semantic | Medium |

## Claude Code Integration

```bash
# Print the embedded skill
exa skill print

# Install skill to ~/.claude/skills/
exa skill add
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | JSON output |
| `--plaintext` | `-p` | Tab-separated output |
| `--no-color` | | Disable colors |
| `--debug` | | Debug logging to stderr |
| `--fields` | | Comma-separated fields for JSON |
| `--jq` | | JQ expression to filter JSON |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `EXA_API_KEY` | API key (required) |
| `EXA_API_URL` | API base URL (default: https://api.exa.ai) |
| `NO_COLOR` | Disable colored output |

## License

MIT
