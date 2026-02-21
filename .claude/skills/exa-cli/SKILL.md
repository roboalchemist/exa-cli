---
name: exa-cli
description: "CLI for the Exa AI search API. Searches the web, finds similar pages, gets AI-powered answers with citations, retrieves page contents, and explores code context. Use when performing web search, content retrieval, finding similar pages, or getting AI answers."
scope: both
---

# exa-cli

CLI for the Exa AI search API. All commands support `--json` for structured output and `--plaintext` for piping.

## When to Use Exa vs Perplexity

| Need | Use |
|------|-----|
| Discover URLs / find sources | **Exa** (`exa search`) |
| Extract content from known URLs | **Exa** (`exa contents`) |
| Find pages similar to a URL | **Exa** (`exa similar`) |
| Code-specific context search | **Exa** (`exa context`) |
| Quick factual Q&A with citations | **Exa** (`exa answer`) or Perplexity `sonar` |
| Synthesized answer with reasoning | **Perplexity** (`sonar-reasoning`) |
| Compare / analyze tradeoffs | **Perplexity** (`sonar-reasoning-pro`) |
| Domain-filtered synthesized answer | **Perplexity** (`-d`) |

**Rule of thumb**: Use **Exa** to find things, use **Perplexity** to understand things.

## Which Command for Which Task

| Task | Command | Notes |
|------|---------|-------|
| Find URLs about a topic | `exa search "query"` | Replaces Tavily search |
| Get full page text from URLs | `exa contents URL` | Replaces Tavily extract; supports `--livecrawl` |
| Quick AI answer with sources | `exa answer "query"` | Grounded Q&A, cheaper than Perplexity |
| Find related pages | `exa similar URL` | Semantic similarity, unique to Exa |
| Search code repos/docs | `exa context "query"` | Code-specific results from Exa Code |
| Check API usage/costs | `exa usage` | Monitor spending |

> **Note**: Exa also has a Research API (async, multi-step) but it's SDK-only — our CLI doesn't support it yet.

<examples>
<example>
Task: Discover URLs about a topic (cost-optimal discovery)

```bash
exa search "hottest AI startups" -n 25 --no-contents --json --fields title,url,score
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
Task: Two-step research — discover then extract (3x cheaper than search with text)

```bash
# Step 1: Find URLs ($0.005)
exa search "python async frameworks" -n 25 --no-contents --json --fields title,url,score

# Step 2: Pull text from top results ($0.003-0.005)
exa contents https://best-result.com https://second-result.com https://third-result.com --json
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
exa similar "https://arxiv.org/abs/2307.06435" -n 25 --exclude-source --json
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
exa search "machine learning" -n 25 --no-contents --category research_paper --start-date 2025-01-01 --json --jq '.results | length'
```

Output:
```
25
```
</example>

<example>
Task: Instant search for real-time low-latency results

```bash
exa search "breaking news AI" -n 25 --type instant --no-contents --json --fields title,url
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

# Search (discovery — no text, max free results)
exa search "query" -n 25 --no-contents --type auto
```

## Commands

| Command | Description |
|---------|-------------|
| `search [query]` | Web search (types: auto, fast, deep, neural, instant) |
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

| Type | Measured Latency | Use |
|------|-----------------|-----|
| `instant` | ~0.3s | Real-time apps needing minimal latency |
| `neural` | ~0.2s | Pure semantic/embeddings search |
| `fast` | ~0.5s | Real-time applications |
| `auto` | ~0.9s | Default, combines methods with reranker |
| `deep` | ~3.8s | Comprehensive, query expansion |

Other commands: `contents` ~0.1s, `similar` ~0.2s, `answer` ~6s.

## Authentication

Priority: `EXA_API_KEY` env var > `~/.exa-auth.json` config file.

## Pricing (Pay-As-You-Go)

| Endpoint | Cost | Notes |
|----------|------|-------|
| `search` (1-25 results) | $5 / 1k ($0.005/req) | **Always use `-n 25`** — same price as `-n 1` |
| `search` (26-100 results) | $25 / 1k ($0.025/req) | 5x more expensive — only if you truly need 26+ |
| `search` (deep) | $15 / 1k ($0.015/req) | |
| `contents` text | $1 / 1k pages ($0.001/pg) | Additive — charged PER RESULT on top of search |
| `contents` highlights | $1 / 1k pages ($0.001/pg) | Additive — separate from text |
| `contents` summary | $1 / 1k pages ($0.001/pg) | Additive — separate from text |
| `similar` | Same tiers as search | **Use `-n 25` here too** (text defaults OFF) |
| `answer` | $5 / 1k ($0.005/req) | Flat rate |
| `context` | ~$0.015/req | Includes highlights |

### Empirical Cost Table (25 results)

| Operation | Cost | Ratio |
|-----------|------|-------|
| `search -n 25 --no-contents` | $0.005 | 1x (baseline) |
| `search -n 25` (text on by default) | $0.030 | 6x |
| `search -n 25 --highlights` | $0.055 | 11x |
| `search -n 25 --summary` | $0.055 | 11x |
| `search -n 25 --highlights --summary` | $0.080 | 16x |
| `similar -n 25` (text off by default) | $0.005 | 1x |
| `similar -n 25 --text` | $0.030 | 6x |
| `contents` (5 URLs, text) | $0.005 | 1x |
| `answer` | $0.005 | 1x |

### Cost-Optimal Research Pattern

**Discovery then selective extraction** saves 3x vs blanket text retrieval:

```bash
# Step 1: Discover URLs — $0.005 (no text)
exa search "query" -n 25 --no-contents --json --fields title,url,score

# Step 2: Pull text from top 3-5 URLs — $0.003-0.005
exa contents URL1 URL2 URL3 --json
```

**Total: ~$0.010** vs **$0.030** for `search -n 25` with default text on all 25 results.

Only use `search` with text enabled (`--text`, the default) when you need to scan text across many results simultaneously.

## For Agents

- **Always use `-n 25`** for `search` and `similar` (same price as 1-25)
- **Prefer two-step search**: `search --no-contents` then `contents` on best URLs (3x cheaper)
- Only use `search` with text (default) when you need to scan all results' text at once
- Always use `--json` for programmatic parsing
- Use `--fields` to reduce output: `--json --fields title,url`
- Use `--jq` for inline filtering: `--json --jq '.results[0].url'`
- Cost is included in JSON responses under `costDollars.total`
- Errors output structured JSON to stderr with `--json`
- Use `exa usage` to check API costs before bulk operations

See [reference/commands.md](reference/commands.md) for complete flag reference.
