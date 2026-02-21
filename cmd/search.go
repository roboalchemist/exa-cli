package cmd

import (
	"fmt"
	"strings"

	"github.com/roboalchemist/exa-cli/pkg/api"
	"github.com/roboalchemist/exa-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	searchNumResults  int
	searchType        string
	searchCategory    string
	searchIncDomains  []string
	searchExcDomains  []string
	searchStartDate   string
	searchEndDate     string
	searchIncludeText string
	searchExcludeText string
	searchText        bool
	searchTextMax     int
	searchHighlights  bool
	searchSummary     bool
	searchNoContents  bool
	searchMaxAge      int
	searchModeration  bool
	searchSubpages    int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the web using Exa AI",
	Long: `Search the web using the Exa AI search API.

Supports multiple search types:
  auto   — Combines methods with reranker (default)
  fast   — <400ms latency, good for real-time
  deep   — Query expansion, comprehensive results
  neural — Pure embeddings-based semantic search

Examples:
  exa search "hottest AI startups"
  exa search "climate change" --type deep -n 20
  exa search "machine learning" --category research_paper
  exa search "golang tutorials" --include-domains go.dev,gobyexample.com
  exa search "AI news" --start-date 2025-01-01 --highlights
  exa search "React hooks" --json --fields title,url,score`,
	Args:       cobra.MinimumNArgs(1),
	SuggestFor: []string{"find", "query", "lookup"},
	RunE:       runSearch,
}

func init() {
	f := searchCmd.Flags()
	f.IntVarP(&searchNumResults, "num-results", "n", 25, "Max results (max 100)")
	f.StringVarP(&searchType, "type", "t", "auto", "Search type: auto|fast|deep|neural")
	f.StringVar(&searchCategory, "category", "", "Category: company|news|research_paper|tweet|github|etc")
	f.StringSliceVar(&searchIncDomains, "include-domains", nil, "Only search these domains")
	f.StringSliceVar(&searchExcDomains, "exclude-domains", nil, "Exclude these domains")
	f.StringVar(&searchStartDate, "start-date", "", "Published after (YYYY-MM-DD)")
	f.StringVar(&searchEndDate, "end-date", "", "Published before (YYYY-MM-DD)")
	f.StringVar(&searchIncludeText, "include-text", "", "Text that must appear in results")
	f.StringVar(&searchExcludeText, "exclude-text", "", "Text that must NOT appear in results")
	f.BoolVar(&searchText, "text", false, "Include full text in results (adds $0.001/result)")
	f.IntVar(&searchTextMax, "text-max-chars", 10000, "Max chars for text content")
	f.BoolVar(&searchHighlights, "highlights", false, "Include LLM-selected highlights")
	f.BoolVar(&searchSummary, "summary", false, "Include LLM summary")
	f.BoolVar(&searchNoContents, "no-contents", false, "Disable all content retrieval")
	f.IntVar(&searchMaxAge, "max-age-hours", -1, "Max cache age (-1=cache, 0=always livecrawl)")
	f.BoolVar(&searchModeration, "moderation", false, "Enable content safety moderation")
	f.IntVar(&searchSubpages, "subpages", 0, "Number of subpages to crawl per result")

	_ = searchCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"auto\tCombines methods with reranker (default)",
			"fast\tLow latency (<400ms)",
			"deep\tQuery expansion, comprehensive",
			"neural\tPure embeddings-based semantic",
		}, cobra.ShellCompDirectiveNoFileComp
	})

	_ = searchCmd.RegisterFlagCompletionFunc("category", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"company", "news", "research_paper", "tweet", "github",
			"linkedin_profile", "pdf", "personal_site",
		}, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	req := &api.SearchRequest{
		Query:      strings.Join(args, " "),
		NumResults: searchNumResults,
	}

	if searchType != "auto" {
		req.Type = searchType
	}
	if searchCategory != "" {
		req.Category = searchCategory
	}
	if len(searchIncDomains) > 0 {
		req.IncludeDomains = searchIncDomains
	}
	if len(searchExcDomains) > 0 {
		req.ExcludeDomains = searchExcDomains
	}
	if searchStartDate != "" {
		req.StartPublishedDate = searchStartDate + "T00:00:00.000Z"
	}
	if searchEndDate != "" {
		req.EndPublishedDate = searchEndDate + "T00:00:00.000Z"
	}
	if searchIncludeText != "" {
		req.IncludeText = searchIncludeText
	}
	if searchExcludeText != "" {
		req.ExcludeText = searchExcludeText
	}
	if searchModeration {
		mod := true
		req.Moderation = &mod
	}

	if !searchNoContents {
		contents := &api.ContentsSpec{}
		hasContents := false

		if searchText {
			contents.Text = &api.TextSpec{MaxCharacters: searchTextMax}
			hasContents = true
		}
		if searchHighlights {
			contents.Highlights = &api.HighlightsSpec{}
			hasContents = true
		}
		if searchSummary {
			contents.Summary = &api.SummarySpec{}
			hasContents = true
		}
		if searchMaxAge >= 0 {
			if searchMaxAge == 0 {
				contents.Livecrawl = "always"
			} else {
				contents.Livecrawl = "fallback"
			}
			hasContents = true
		}
		if searchSubpages > 0 {
			contents.Subpages = searchSubpages
			hasContents = true
		}

		if hasContents {
			req.Contents = contents
		}
	}

	resp, err := client.Search(newContext(), req)
	if err != nil {
		return err
	}

	opts := GetOutputOptions()

	if opts.Mode == output.ModeJSON {
		return output.RenderJSON(resp, opts)
	}

	td := output.TableData{
		Headers: []string{"TITLE", "URL", "DATE", "SCORE"},
	}

	for _, r := range resp.Results {
		date := ""
		if r.PublishedDate != "" {
			if len(r.PublishedDate) >= 10 {
				date = r.PublishedDate[:10]
			} else {
				date = r.PublishedDate
			}
		}
		title := truncateStr(r.Title, 50)
		td.Rows = append(td.Rows, []string{title, r.URL, date, fmt.Sprintf("%.2f", r.Score)})
	}

	footer := fmt.Sprintf("%d results | Type: %s", len(resp.Results), searchType)
	if resp.CostDollars != nil {
		footer = fmt.Sprintf("Cost: $%.4f | %s", resp.CostDollars.Total, footer)
	}
	td.Footer = footer

	return output.RenderTable(td, resp, opts)
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
