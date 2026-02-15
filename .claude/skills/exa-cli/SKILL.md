---
name: exa-cli
description: "CLI for the Exa AI search API. Searches the web, finds similar pages, gets AI-powered answers with citations, retrieves page contents, and explores code context. Use when performing web search, content retrieval, finding similar pages, or getting AI answers."
---

# exa-cli

CLI for the Exa AI search API. All commands support `--json` for structured output and `--plaintext` for piping.

<examples>
<example>
Task: Search for AI startups and get JSON with specific fields

```bash
exa search "hottest AI startups" -n 5 --json --fields title,url,score
```

Output:
```json
[
  {"title": "Adept: Useful General Intelligence", "url": "https://www.adept.ai/", "score": 0.95},
  {"title": "Home | Tenyx, Inc.", "url": "https://www.tenyx.com/", "score": 0.89}
]
```
</example>

<example>
Task: Get an AI answer with citations

```bash
exa answer "What is the capital of France?" --json
```

Output:
```json
{
  "answer": "The capital of France is Paris...",
  "citations": [{"title": "...", "url": "https://..."}],
  "costDollars": {"total": 0.005}
}
```
</example>

<example>
Task: Find pages similar to a research paper

```bash
exa similar "https://arxiv.org/abs/2307.06435" -n 5 --exclude-source --json
```
</example>

<example>
Task: Get page contents with summary

```bash
exa contents https://example.com --summary --json --fields title,url,summary
```
</example>

<example>
Task: Search with filters and count results using jq

```bash
exa search "machine learning" -n 10 --category research_paper --start-date 2025-01-01 --json --jq '.results | length'
```

Output:
```
10
```
</example>
</examples>

## Quick Start

```bash
# Install
brew install roboalchemist/tap/exa-cli

# Authenticate
export EXA_API_KEY="your-key-here"
# Or: exa auth

# Search
exa search "query" -n 10 --type auto
```

## Commands

| Command | Description |
|---------|-------------|
| `search [query]` | Web search (types: auto, fast, deep, neural) |
| `answer [query]` | AI answer with citations (supports `--stream`) |
| `similar [url]` | Find semantically similar pages |
| `contents [urls...]` | Retrieve page text, highlights, summaries |
| `context [query]` | Code context from Exa Code |
| `usage` | API usage stats and costs |
| `auth` | Configure API key |

## Output Formats

- `--json` (`-j`) — Structured JSON (all commands)
- `--plaintext` (`-p`) — Tab-separated for piping
- `--fields title,url,score` — Filter JSON fields
- `--jq '.results[] | .url'` — Built-in JQ filtering

## Search Types

| Type | Latency | Use |
|------|---------|-----|
| `auto` | Medium | Default, combines methods with reranker |
| `fast` | <400ms | Real-time applications |
| `deep` | Higher | Comprehensive, query expansion |
| `neural` | Medium | Pure semantic/embeddings search |

## Authentication

Priority: `EXA_API_KEY` env var > `~/.exa-auth.json` config file.

## For Agents

- Always use `--json` for programmatic parsing
- Use `--fields` to reduce output: `--json --fields title,url`
- Use `--jq` for inline filtering: `--json --jq '.results[0].url'`
- Cost is included in JSON responses under `costDollars.total`
- Errors output structured JSON to stderr with `--json`

See [reference/commands.md](reference/commands.md) for complete flag reference.
