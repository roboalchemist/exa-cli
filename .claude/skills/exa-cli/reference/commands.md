# exa-cli Command Reference

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | JSON output |
| `--plaintext` | `-p` | Tab-separated output for piping |
| `--no-color` | | Disable colored output |
| `--debug` | | Verbose logging to stderr |
| `--fields` | | Comma-separated fields for JSON output |
| `--jq` | | JQ expression to filter JSON output |

## `exa search [query]`

Search the web using Exa AI.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--num-results` | `-n` | 25 | Max results (max 100). 1-25 same price tier. |
| `--type` | `-t` | auto | Search type: auto\|fast\|deep\|neural |
| `--category` | | | company\|news\|research_paper\|tweet\|github\|linkedin_profile\|pdf\|personal_site |
| `--include-domains` | | | Only search these domains (comma-separated) |
| `--exclude-domains` | | | Exclude these domains |
| `--start-date` | | | Published after (YYYY-MM-DD) |
| `--end-date` | | | Published before (YYYY-MM-DD) |
| `--include-text` | | | Text that must appear in results |
| `--exclude-text` | | | Text that must NOT appear |
| `--text` | | false | Include full text in results (adds $0.001/result) |
| `--text-max-chars` | | 10000 | Max chars for text content |
| `--highlights` | | false | Include LLM-selected highlights |
| `--summary` | | false | Include LLM summary |
| `--no-contents` | | false | Disable all content retrieval |
| `--max-age-hours` | | -1 | Max cache age (-1=cache, 0=always livecrawl) |
| `--moderation` | | false | Enable content safety moderation |
| `--subpages` | | 0 | Number of subpages to crawl per result |

## `exa answer [query]`

Get an AI-powered answer with citations.

| Flag | Default | Description |
|------|---------|-------------|
| `--stream` | false | Stream the answer token by token |
| `--text` | false | Include full text in citation sources |
| `--output-schema` | | Path to JSON schema file for structured output |

## `exa similar [url]`

Find pages similar to a URL.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--num-results` | `-n` | 25 | Max results. 1-25 same price tier. |
| `--exclude-source` | | false | Exclude the source domain |
| `--include-domains` | | | Only include these domains |
| `--exclude-domains` | | | Exclude these domains |
| `--start-date` | | | Published after (YYYY-MM-DD) |
| `--end-date` | | | Published before (YYYY-MM-DD) |
| `--text` | | false | Include full text |
| `--highlights` | | false | Include highlights |
| `--category` | | | Category filter |

## `exa contents [urls...]`

Get page contents by URL.

| Flag | Default | Description |
|------|---------|-------------|
| `--text` | true | Include full text |
| `--text-max-chars` | 10000 | Max chars for text |
| `--highlights` | false | Include highlights |
| `--summary` | false | Include summary |
| `--max-age-hours` | -1 | Content freshness (-1=cache, 0=always livecrawl) |
| `--subpages` | 0 | Subpages to crawl |

## `exa context [query]`

Get code context from Exa Code.

| Flag | Default | Description |
|------|---------|-------------|
| `--tokens` | 0 | Token limit (0=dynamic) |

## `exa usage`

Show API usage and costs.

| Flag | Default | Description |
|------|---------|-------------|
| `--start-date` | 30 days ago | Start of period |
| `--end-date` | now | End of period |
| `--key-id` | | Specific API key ID (auto-detects if omitted) |

## `exa auth`

Configure API key interactively. Stores in `~/.exa-auth.json` (mode 0600).

Non-interactive mode requires `EXA_API_KEY` env var.

## `exa completion [bash|zsh|fish|powershell]`

Generate shell completion scripts.

## `exa skill print|add`

Print or install the embedded Claude Code skill.
