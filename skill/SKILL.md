# exa-cli

CLI for the Exa AI search API. Searches the web, finds similar pages, gets AI-powered answers, retrieves page contents, and explores code context. Use when performing web search, finding similar content, getting AI answers with citations, or retrieving page text.

<examples>
<example>
Task: Search for AI startups
```bash
exa search "hottest AI startups" -n 5
exa search "hottest AI startups" -n 5 --json
```
</example>

<example>
Task: Find research papers similar to one
```bash
exa similar "https://arxiv.org/abs/2307.06435" -n 10
```
</example>

<example>
Task: Get an AI answer with citations
```bash
exa answer "What is the current valuation of SpaceX?"
exa answer "How does CRISPR work?" --stream
```
</example>

<example>
Task: Get page contents
```bash
exa contents https://example.com --text --summary
exa contents https://example.com --json
```
</example>

<example>
Task: Search for code context
```bash
exa context "React hooks state management" --tokens 5000
```
</example>

<example>
Task: Check API usage
```bash
exa usage
exa usage --start-date 2025-01-01 --json
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

## Output Formats

- `--json` (`-j`) — Structured JSON (all commands)
- `--plaintext` (`-p`) — Tab-separated for piping
- `--fields` — Filter JSON fields: `--json --fields title,url,score`
- `--jq` — JQ expression: `--json --jq '.results[] | .url'`

## Authentication

Priority: `EXA_API_KEY` env var > `~/.exa-auth.json` config file.

See [reference/commands.md](reference/commands.md) for complete command reference.
